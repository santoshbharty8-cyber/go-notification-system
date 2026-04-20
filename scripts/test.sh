#!/bin/bash

set -e

echo "🧪 Running tests with coverage..."

PKGS=$(go list ./... | grep -v "internal/bootstrap" \
                         | grep -v "internal/logger" \
                         | grep -v "internal/redisclient" \
                         | grep -v "tests/helpers" \
                         | grep -v "/cmd/")

echo "📦 Testing packages:"
echo "$PKGS"

go test $PKGS -coverprofile=coverage.out -covermode=atomic

echo ""
echo "📊 Coverage Summary:"
SUMMARY=$(go tool cover -func=coverage.out | grep total)
echo "$SUMMARY"

COVERAGE=$(echo "$SUMMARY" | awk '{print $3}' | sed 's/%//')

echo "✅ Total Coverage: $COVERAGE%"

THRESHOLD=90
COVERAGE_INT=${COVERAGE%.*}

if [ "$COVERAGE_INT" -lt "$THRESHOLD" ]; then
  echo "❌ Coverage below threshold ($THRESHOLD%)"
  exit 1
fi

echo "🎉 Coverage check passed!"