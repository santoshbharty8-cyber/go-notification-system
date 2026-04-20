#!/bin/bash

set -e

echo "🚀 Running full quality checks..."

# =========================
# 1. LINT
# =========================
echo ""
echo "🧹 Running lint..."
golangci-lint run

# =========================
# 2. SECURITY SCAN (GOSEC)
# =========================
echo ""
echo "🔒 Running gosec..."
gosec ./...

# =========================
# 3. VULNERABILITY SCAN
# =========================
echo ""
echo "🛡️ Running govulncheck..."

if command -v govulncheck >/dev/null 2>&1; then
  govulncheck ./...
else
  echo "⚠️ govulncheck not found, running via go run..."
  go run golang.org/x/vuln/cmd/govulncheck@latest ./...
fi

# =========================
# 4. TEST + COVERAGE
# =========================
echo ""
echo "🧪 Running tests with coverage..."

# 🔥 Exclude infra packages
PKGS=$(go list ./... | grep -v "internal/bootstrap" \
                         | grep -v "internal/logger" \
                         | grep -v "internal/redisclient" \
                         | grep -v "tests/helpers" \
                         | grep -v "/cmd/")

echo "📦 Testing packages:"
echo "$PKGS"

go test $PKGS -coverprofile=coverage.out -covermode=atomic

# =========================
# 5. COVERAGE SUMMARY
# =========================
echo ""
echo "📊 Coverage Summary:"
SUMMARY=$(go tool cover -func=coverage.out | grep total)
echo "$SUMMARY"

COVERAGE=$(echo "$SUMMARY" | awk '{print $3}' | sed 's/%//')

echo "✅ Total Coverage: $COVERAGE%"

# =========================
# 6. COVERAGE CHECK
# =========================
THRESHOLD=90
COVERAGE_INT=${COVERAGE%.*}

if [ "$COVERAGE_INT" -lt "$THRESHOLD" ]; then
  echo "❌ Coverage below threshold ($THRESHOLD%)"
  exit 1
fi

echo "🎉 Coverage check passed!"

# =========================
# 7. OPTIONAL HTML REPORT
# =========================
if [ -f coverage.out ]; then
  echo ""
  echo "📄 Generating coverage report (coverage.html)..."
  go tool cover -html=coverage.out -o coverage.html
fi

echo ""
echo "✅ ALL CHECKS PASSED (Lint + Security + Tests + Coverage)"