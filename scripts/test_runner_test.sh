#!/usr/bin/env bash
# Test execution verification for scripts/run-all-tests.sh

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
RUNNER="$SCRIPT_DIR/run-all-tests.sh"

echo "Running CLI tests for $RUNNER..."

# Test 1: Help flag
"$RUNNER" --help > /dev/null
echo "✓ --help test passed."

# Test 2: Invalid flag handling
set +e
"$RUNNER" --invalid-flag > /dev/null 2>&1
EXIT_CODE=$?
set -e

if [ $EXIT_CODE -ne 1 ]; then
    echo "✗ Expected exit code 1 for invalid flag, got $EXIT_CODE" >&2
    exit 1
fi
echo "✓ Invalid flag rejection test passed."

# Test 3: Backend tier execution
"$RUNNER" --backend > /dev/null
echo "✓ --backend flag test passed."

# Test 4: Frontend tier execution
"$RUNNER" --frontend > /dev/null
echo "✓ --frontend flag test passed."

echo "All test runner CLI assertions passed!"

