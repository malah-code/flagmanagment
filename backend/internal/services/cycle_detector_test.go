package services

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCycleDetectorService(t *testing.T) {
	detector := NewCycleDetectorService()
	ctx := context.Background()

	t.Run("No Cycle - Nil Parent", func(t *testing.T) {
		targetID := uuid.New()
		err := detector.DetectCycle(ctx, targetID, nil, map[uuid.UUID]*uuid.UUID{})
		assert.NoError(t, err)
	})

	t.Run("Immediate Self-Reference", func(t *testing.T) {
		targetID := uuid.New()
		err := detector.DetectCycle(ctx, targetID, &targetID, map[uuid.UUID]*uuid.UUID{})
		assert.ErrorIs(t, err, ErrCircularDependency)
	})

	t.Run("A -> B -> A Cycle", func(t *testing.T) {
		targetID := uuid.New()
		parentID := uuid.New()
		
		existing := map[uuid.UUID]*uuid.UUID{
			parentID: &targetID,
		}

		err := detector.DetectCycle(ctx, targetID, &parentID, existing)
		assert.ErrorIs(t, err, ErrCircularDependency)
	})

	t.Run("Max Depth Exceeded", func(t *testing.T) {
		targetID := uuid.New()
		id1 := uuid.New()
		id2 := uuid.New()
		id3 := uuid.New()
		id4 := uuid.New()

		existing := map[uuid.UUID]*uuid.UUID{
			id1: &id2,
			id2: &id3,
			id3: &id4,
		}

		err := detector.DetectCycle(ctx, targetID, &id1, existing)
		assert.ErrorIs(t, err, ErrMaxDependencyDepth)
	})

	t.Run("Valid Chain", func(t *testing.T) {
		targetID := uuid.New()
		id1 := uuid.New()
		id2 := uuid.New()

		existing := map[uuid.UUID]*uuid.UUID{
			id1: &id2,
			id2: nil,
		}

		err := detector.DetectCycle(ctx, targetID, &id1, existing)
		assert.NoError(t, err)
	})
}
