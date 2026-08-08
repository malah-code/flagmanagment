package upstream

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/flagmanagment/proxy/internal/broadcaster"
	"github.com/flagmanagment/proxy/internal/config"
	"github.com/flagmanagment/proxy/internal/store"
	pb "github.com/flagmanagment/proxy/pkg/gen/sdk/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type UpstreamClient struct {
	cfg         *config.Config
	store       *store.RulesetStore
	broadcaster *broadcaster.Broadcaster
	state       *UpstreamState
}

func NewUpstreamClient(cfg *config.Config, store *store.RulesetStore, broadcaster *broadcaster.Broadcaster, state *UpstreamState) *UpstreamClient {
	return &UpstreamClient{
		cfg:         cfg,
		store:       store,
		broadcaster: broadcaster,
		state:       state,
	}
}

func (c *UpstreamClient) State() *UpstreamState {
	return c.state
}

func (c *UpstreamClient) Run(ctx context.Context) {
	backoff := 1 * time.Second
	maxBackoff := 60 * time.Second

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Upstream client stopping")
			return
		default:
		}

		start := time.Now()
		err := c.connectAndStream(ctx)
		c.state.SetConnected(false)

		if time.Since(start) > 1*time.Minute {
			backoff = 1 * time.Second
		}

		if ctx.Err() != nil {
			return
		}

		// Calculate jittered backoff: ±25%
		jitter := time.Duration((rand.Float64()*0.5 - 0.25) * float64(backoff))
		waitDuration := backoff + jitter

		log.Warn().Err(err).Dur("retry_in", waitDuration).Msg("Upstream stream disconnected, retrying")

		select {
		case <-ctx.Done():
			return
		case <-time.After(waitDuration):
		}

		backoff *= 2
		if backoff > maxBackoff {
			backoff = maxBackoff
		}
	}
}

func (c *UpstreamClient) connectAndStream(ctx context.Context) error {
	var opts []grpc.DialOption
	if c.cfg.UpstreamTLS {
		creds := credentials.NewClientTLSFromCert(nil, "")
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	log.Info().Str("backend_addr", c.cfg.BackendAddr).Msg("Connecting to FlagManagment backend")
	conn, err := grpc.DialContext(ctx, c.cfg.BackendAddr, opts...)
	if err != nil {
		return fmt.Errorf("failed to dial backend: %w", err)
	}
	defer conn.Close()

	client := pb.NewSDKServiceClient(conn)

	// Auth metadata context
	md := metadata.Pairs("authorization", "Bearer "+c.cfg.EnvToken)
	authCtx := metadata.NewOutgoingContext(ctx, md)

	// 1. Fetch Snapshot (Bootstrap)
	snapshot, err := client.FetchSnapshot(authCtx, &pb.SnapshotRequest{
		EnvironmentToken: c.cfg.EnvToken,
	})
	if err != nil {
		return fmt.Errorf("failed to fetch initial snapshot: %w", err)
	}

	c.store.Set(snapshot)
	c.state.SetConnected(true)
	log.Info().Str("version", snapshot.Version).Int("flags_count", len(snapshot.Flags)).Msg("Successfully bootstrapped ruleset snapshot from backend")

	// Notify connected downstream clients of full reset / new version
	c.broadcaster.Broadcast(&pb.RulesetDelta{
		FullReset: true,
		Version:   snapshot.Version,
	})

	// 2. Open Stream
	stream, err := client.StreamRulesets(authCtx, &pb.StreamRequest{
		EnvironmentToken: c.cfg.EnvToken,
	})
	if err != nil {
		return fmt.Errorf("failed to establish ruleset stream: %w", err)
	}

	for {
		delta, err := stream.Recv()
		if err == io.EOF {
			return fmt.Errorf("stream closed by server")
		}
		if err != nil {
			return fmt.Errorf("stream error: %w", err)
		}

		c.state.RecordDelta()
		log.Debug().Str("version", delta.Version).Bool("full_reset", delta.FullReset).Msg("Received delta update from upstream")

		if delta.FullReset {
			// Refetch full snapshot if requested
			freshSnapshot, err := client.FetchSnapshot(authCtx, &pb.SnapshotRequest{
				EnvironmentToken: c.cfg.EnvToken,
			})
			if err == nil {
				c.store.Set(freshSnapshot)
			}
		}

		// Fan out to downstream SDK clients
		c.broadcaster.Broadcast(delta)
	}
}
