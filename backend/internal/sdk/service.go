package sdk

import (
	"context"
	"github.com/flagmanagment/backend/internal/cache"
	"github.com/flagmanagment/backend/internal/middleware"
	"github.com/flagmanagment/backend/internal/repository"
	pb "github.com/flagmanagment/backend/pkg/gen/sdk/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedSDKServiceServer
	store       repository.Store
	cacheClient *cache.Client
}

func NewServer(store repository.Store, cacheClient *cache.Client) *Server {
	return &Server{
		store:       store,
		cacheClient: cacheClient,
	}
}

func (s *Server) FetchSnapshot(ctx context.Context, req *pb.SnapshotRequest) (*pb.RulesetSnapshot, error) {
	env := middleware.GetEnvironmentFromContext(ctx)
	if env == nil {
		// Fallback to token inside request if metadata interceptor wasn't used
		var err error
		env, err = middleware.ValidateSDKToken(ctx, s.store, req.GetEnvironmentToken())
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid environment token")
		}
	}

	// Attempt to retrieve from Redis cache
	snapshot, err := s.cacheClient.GetRulesetSnapshot(ctx, env.ID.String())
	if err != nil || snapshot == nil {
		// Per clarification: Fail fast with 503 Service Unavailable if cache is missed/unavailable
		return nil, status.Error(codes.Unavailable, "Redis cache unavailable or ruleset not populated")
	}

	pbFlags := make([]*pb.FlagRule, len(snapshot.Flags))
	for i, f := range snapshot.Flags {
		pbFlags[i] = &pb.FlagRule{
			Key:                f.Key,
			Type:               f.Type,
			Enabled:            f.Enabled,
			DefaultVariation:   f.DefaultVariation,
			TargetingRulesJson: f.TargetingRules,
		}
	}

	return &pb.RulesetSnapshot{
		Version: snapshot.Version,
		Flags:   pbFlags,
	}, nil
}

func (s *Server) StreamRulesets(req *pb.StreamRequest, stream pb.SDKService_StreamRulesetsServer) error {
	ctx := stream.Context()
	env := middleware.GetEnvironmentFromContext(ctx)
	if env == nil {
		var err error
		env, err = middleware.ValidateSDKToken(ctx, s.store, req.GetEnvironmentToken())
		if err != nil {
			return status.Error(codes.Unauthenticated, "invalid environment token")
		}
	}

	pubsub := s.cacheClient.SubscribeRulesetUpdates(ctx, env.ID.String())
	defer pubsub.Close()

	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-ch:
			if !ok {
				return nil
			}
			// Send lightweight delta notification
			delta := &pb.RulesetDelta{
				FullReset: false,
				Version:   msg.Payload,
			}
			if err := stream.Send(delta); err != nil {
				return err
			}
		}
	}
}
