package upstream

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/flagmanagment/proxy/internal/broadcaster"
	"github.com/flagmanagment/proxy/internal/config"
	"github.com/flagmanagment/proxy/internal/store"
	pb "github.com/flagmanagment/proxy/pkg/gen/sdk/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type mockBackendSDKServer struct {
	pb.UnimplementedSDKServiceServer
	snapshot *pb.RulesetSnapshot
}

func (m *mockBackendSDKServer) FetchSnapshot(ctx context.Context, req *pb.SnapshotRequest) (*pb.RulesetSnapshot, error) {
	return m.snapshot, nil
}

func (m *mockBackendSDKServer) StreamRulesets(req *pb.StreamRequest, stream pb.SDKService_StreamRulesetsServer) error {
	// Send one delta and exit
	_ = stream.Send(&pb.RulesetDelta{Version: "v1.0.1"})
	return nil
}

func TestUpstreamClient_BootstrapAndStream(t *testing.T) {
	// Start in-memory mock gRPC server
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err)
	defer lis.Close()

	grpcServer := grpc.NewServer()
	mockServer := &mockBackendSDKServer{
		snapshot: &pb.RulesetSnapshot{
			Version: "v1.0.0",
			Flags: []*pb.FlagRule{
				{Key: "test-flag", Enabled: true},
			},
		},
	}
	pb.RegisterSDKServiceServer(grpcServer, mockServer)

	go func() {
		_ = grpcServer.Serve(lis)
	}()
	defer grpcServer.Stop()

	cfg := &config.Config{
		BackendAddr: lis.Addr().String(),
		EnvToken:    "env_test_token",
		UpstreamTLS: false,
	}

	rStore := store.NewRulesetStore()
	bCaster := broadcaster.NewBroadcaster()
	state := NewUpstreamState()

	client := NewUpstreamClient(cfg, rStore, bCaster, state)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go client.Run(ctx)

	// Wait for store to be populated by bootstrap
	assert.Eventually(t, func() bool {
		return rStore.Version() == "v1.0.0"
	}, 1*time.Second, 50*time.Millisecond)

	snap := rStore.Get()
	assert.NotNil(t, snap)
	assert.Equal(t, 1, len(snap.Flags))
	assert.Equal(t, "test-flag", snap.Flags[0].Key)
}
