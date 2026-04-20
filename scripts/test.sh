#!/bin/bash

set -e

echo "🚀 Running tests with coverage..."

# 🔥 Exclude infra packages PROPERLY
PKGS=$(go list ./... | grep -v "internal/bootstrap" \
                         | grep -v "internal/logger" \
                         | grep -v "internal/redisclient" \
                         | grep -v "tests/helpers" \
                         | grep -v "/cmd/")

echo "📦 Testing packages:"
echo "$PKGS"

# Run tests
go test $PKGS -coverprofile=coverage.out -covermode=atomic

# Show summary
echo ""
echo "📊 Coverage Summary:"
SUMMARY=$(go tool cover -func=coverage.out | grep total)
echo "$SUMMARY"

# Extract number
COVERAGE=$(echo "$SUMMARY" | awk '{print $3}' | sed 's/%//')

echo "✅ Total Coverage: $COVERAGE%"

# 🔥 Threshold check WITHOUT bc (portable)
THRESHOLD=85
COVERAGE_INT=${COVERAGE%.*}

if [ "$COVERAGE_INT" -lt "$THRESHOLD" ]; then
  echo "❌ Coverage below threshold ($THRESHOLD%)"
  exit 1
fi

echo "🎉 Coverage check passed!"