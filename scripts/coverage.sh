#!/usr/bin/env bash
set -euo pipefail

echo "Running tests and generating raw coverage..."
go test ./... -coverprofile=coverage.raw.out

echo "Filtering coverage (excluding generated & entry code)..."

{
  head -n 1 coverage.raw.out
  grep -vE '(/cmd/|/ent/|wire\.go|wire_gen\.go|\.pb\.go|mock_|/vendor/|/main\.go)' coverage.raw.out | tail -n +2
} > coverage.out

echo "Calculating final coverage..."
COVERAGE=$(go tool cover -func=coverage.out | grep '^total:' | awk '{print $3}' | sed 's/%//')

if [[ -z "$COVERAGE" ]]; then
  echo "❌ ERROR: Coverage could not be calculated"
  exit 1
fi

echo "Final coverage: $COVERAGE%"

if (( $(echo "$COVERAGE < 80" | bc -l) )); then
  echo "❌ Coverage gate failed (minimum 80%)"
  exit 1
fi

echo "Coverage gate passed"
