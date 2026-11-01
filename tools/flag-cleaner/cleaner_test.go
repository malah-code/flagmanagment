package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCleanerScanAndDryRun(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "flagcleaner_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFilePath := filepath.Join(tempDir, "sample.ts")
	code := `
import { useFlag } from '@flagmanagment/react-sdk';

export function CheckoutButton() {
  const isV2 = useFlag('checkout-v2', false);
  if (isV2) {
    return <NewCheckout />;
  }
  return <LegacyCheckout />;
}
`
	if err := os.WriteFile(testFilePath, []byte(code), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	cleaner := NewCleaner()

	// Test Scan
	results, err := cleaner.Scan(tempDir, "checkout-v2")
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].FlagKey != "checkout-v2" {
		t.Errorf("expected flag key 'checkout-v2', got %q", results[0].FlagKey)
	}

	// Test Clean (Dry Run)
	res, err := cleaner.CleanFlag(testFilePath, "checkout-v2", true, true)
	if err != nil {
		t.Fatalf("clean failed: %v", err)
	}
	if res.ReferencesRemoved != 1 {
		t.Errorf("expected 1 reference removed, got %d", res.ReferencesRemoved)
	}
	if !strings.Contains(res.Diff, "checkout-v2") {
		t.Errorf("expected diff to mention 'checkout-v2', got: %s", res.Diff)
	}

	// Verify original file remained unchanged in dry run
	unchanged, _ := os.ReadFile(testFilePath)
	if !strings.Contains(string(unchanged), "checkout-v2") {
		t.Errorf("dry run should not have modified file on disk")
	}

	// Test Clean with Write
	resWrite, err := cleaner.CleanFlag(testFilePath, "checkout-v2", true, false)
	if err != nil {
		t.Fatalf("clean with write failed: %v", err)
	}
	if resWrite.ReferencesRemoved != 1 {
		t.Errorf("expected 1 reference removed, got %d", resWrite.ReferencesRemoved)
	}

	modified, _ := os.ReadFile(testFilePath)
	if strings.Contains(string(modified), "checkout-v2") {
		t.Errorf("expected 'checkout-v2' to be removed from file after write")
	}
}
