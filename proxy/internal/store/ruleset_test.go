package store

import (
	"sync"
	"testing"

	pb "github.com/flagmanagment/proxy/pkg/gen/sdk/v1"
	"github.com/stretchr/testify/assert"
)

func TestRulesetStore_GetSetVersion(t *testing.T) {
	s := NewRulesetStore()

	assert.Nil(t, s.Get())
	assert.Equal(t, "", s.Version())

	snapshot := &pb.RulesetSnapshot{
		Version: "v1.0.0",
		Flags: []*pb.FlagRule{
			{Key: "feature-a", Enabled: true},
		},
	}

	s.Set(snapshot)

	res := s.Get()
	assert.NotNil(t, res)
	assert.Equal(t, "v1.0.0", s.Version())
	assert.Equal(t, 1, len(res.Flags))
	assert.Equal(t, "feature-a", res.Flags[0].Key)
}

func TestRulesetStore_Concurrency(t *testing.T) {
	s := NewRulesetStore()
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func(v int) {
			defer wg.Done()
			s.Set(&pb.RulesetSnapshot{Version: "v1"})
		}(i)

		go func() {
			defer wg.Done()
			_ = s.Get()
			_ = s.Version()
		}()
	}

	wg.Wait()
}
