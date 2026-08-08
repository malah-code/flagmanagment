package server

import (
	"context"
	"strings"

	"github.com/flagmanagment/proxy/internal/broadcaster"
	"github.com/flagmanagment/proxy/internal/config"
	"github.com/flagmanagment/proxy/internal/store"
	pb "github.com/flagmanagment/proxy/pkg/gen/sdk/v1"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ProxyServer struct {
	pb.UnimplementedSDKServiceServer
	cfg         *config.Config
	store       *store.RulesetStore
	broadcaster *broadcaster.Broadcaster
}

func NewProxyServer(cfg *config.Config, store *store.RulesetStore, broadcaster *broadcaster.Broadcaster) *ProxyServer {
	return &ProxyServer{
		cfg:         cfg,
		store:       store,
		broadcaster: broadcaster,
	}
}

func (s *ProxyServer) validateToken(ctx context.Context, tokenInReq string) error {
	token := tokenInReq
	if token == "" {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			authHeader := md.Get("authorization")
			if len(authHeader) > 0 {
				token = strings.TrimPrefix(authHeader[0], "Bearer ")
			}
		}
	}

	if s.cfg.EnvToken != "" && token != s.cfg.EnvToken {
		return status.Error(codes.Unauthenticated, "invalid environment token")
	}
	return nil
}

func (s *ProxyServer) FetchSnapshot(ctx context.Context, req *pb.SnapshotRequest) (*pb.RulesetSnapshot, error) {
	if err := s.validateToken(ctx, req.GetEnvironmentToken()); err != nil {
		return nil, err
	}

	snapshot := s.store.Get()
	if snapshot == nil {
		return nil, status.Error(codes.Unavailable, "proxy starting up or ruleset unavailable")
	}

	return snapshot, nil
}

func (s *ProxyServer) StreamRulesets(req *pb.StreamRequest, stream pb.SDKService_StreamRulesetsServer) error {
	ctx := stream.Context()
	if err := s.validateToken(ctx, req.GetEnvironmentToken()); err != nil {
		return err
	}

	clientID := uuid.New().String()
	clientChan := make(chan *pb.RulesetDelta, 10)
	client := &broadcaster.Client{
		ID: clientID,
		Ch: clientChan,
	}

	s.broadcaster.Register(client)
	defer s.broadcaster.Deregister(clientID)

	log.Debug().Str("client_id", clientID).Msg("Downstream SDK client connected")

	// Immediately send initial delta to notify current version
	if version := s.store.Version(); version != "" {
		_ = stream.Send(&pb.RulesetDelta{
			FullReset: true,
			Version:   version,
		})
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case delta, ok := <-clientChan:
			if !ok {
				return nil
			}
			if err := stream.Send(delta); err != nil {
				return err
			}
		}
	}
}
