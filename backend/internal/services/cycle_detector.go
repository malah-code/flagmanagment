package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrCircularDependency = errors.New("circular dependency detected")
	ErrMaxDependencyDepth = errors.New("maximum dependency depth exceeded")
)

const MaxDependencyDepth = 3

// CycleDetectorService provides logic for detecting cycles in flag dependencies.
type CycleDetectorService struct {
}

// NewCycleDetectorService creates a new CycleDetectorService.
func NewCycleDetectorService() *CycleDetectorService {
	return &CycleDetectorService{}
}

// DetectCycle checks if setting `proposedParentID` as the parent of `targetFlagID`
// would create a cycle, or exceed the maximum allowed depth.
// existingFlags is a map of flag ID to its current parent ID.
func (s *CycleDetectorService) DetectCycle(ctx context.Context, targetFlagID uuid.UUID, proposedParentID *uuid.UUID, existingFlags map[uuid.UUID]*uuid.UUID) error {
	if proposedParentID == nil {
		return nil
	}

	if *proposedParentID == targetFlagID {
		return ErrCircularDependency
	}

	depth := 1
	currentParentID := proposedParentID

	for currentParentID != nil {
		if depth > MaxDependencyDepth {
			return ErrMaxDependencyDepth
		}

		nextParentPtr, exists := existingFlags[*currentParentID]
		if !exists || nextParentPtr == nil {
			break
		}
		
		if *nextParentPtr == targetFlagID {
			return ErrCircularDependency
		}

		currentParentID = nextParentPtr
		depth++
	}

	return nil
}
