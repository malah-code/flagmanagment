package utils

import (
	"fmt"
	"github.com/spaolacci/murmur3"
)

// BucketUser hashes the targeting key and returns a bucket between 0 and 99.
func BucketUser(targetingKey string) int {
	// Use MurmurHash3 32-bit x86 variant.
	hasher := murmur3.New32()
	_, _ = hasher.Write([]byte(targetingKey))
	hashValue := hasher.Sum32()
	
	// Convert to a bucket (0-99).
	return int(hashValue % 100)
}

// GetVariant returns the variant based on the user's bucket and rollout percentages.
// rolloutMap is a map of VariantName -> Percentage (0-100).
func GetVariant(targetingKey string, rolloutMap map[string]int) (string, error) {
	bucket := BucketUser(targetingKey)
	
	currentThreshold := 0
	for variant, percentage := range rolloutMap {
		currentThreshold += percentage
		if bucket < currentThreshold {
			return variant, nil
		}
	}
	
	return "", fmt.Errorf("no variant matched, bucket %d exceeded total percentage %d", bucket, currentThreshold)
}
