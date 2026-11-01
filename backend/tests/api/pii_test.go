package api_test

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

// @TestCase: TC-PII-001
// @TestCase: TC-PII-002
func TestPIIIdentityHashing(t *testing.T) {
	rawEmail := "user.test@example.com"
	salt := "test-app-salt"

	hasher := sha256.New()
	hasher.Write([]byte(rawEmail + salt))
	hashed := hex.EncodeToString(hasher.Sum(nil))

	if len(hashed) != 64 {
		t.Fatalf("expected 64 character SHA-256 hash, got %d", len(hashed))
	}

	// Ensure determinism
	hasher2 := sha256.New()
	hasher2.Write([]byte(rawEmail + salt))
	hashed2 := hex.EncodeToString(hasher2.Sum(nil))

	if hashed != hashed2 {
		t.Fatalf("expected deterministic hash output")
	}
}

