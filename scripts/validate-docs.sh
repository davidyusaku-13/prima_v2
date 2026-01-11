#!/bin/bash
# PRIMA Documentation Validation Script
# Run this in CI or before commits to catch documentation drift

set -e

DOCS_DIR="docs"
BACKEND_DIR="backend"
FRONTEND_DIR="frontend"

echo "=== PRIMA Documentation Validator ==="
echo ""

ERRORS=0

# 1. Check Go version matches documentation
echo "Checking Go version..."
BACKEND_GO_VERSION=$(grep "^go " "$BACKEND_DIR/go.mod" | awk '{print $2}')
DOCS_GO_VERSION=$(grep -oP "Go \K[\d.]+" "$DOCS_DIR/technology-stack.md" | head -1)

if [ "$BACKEND_GO_VERSION" != "$DOCS_GO_VERSION" ]; then
    echo "  ERROR: go.mod says Go $BACKEND_GO_VERSION but docs say $DOCS_GO_VERSION"
    ((ERRORS++))
else
    echo "  OK: Go version $BACKEND_GO_VERSION"
fi

# 2. Check Svelte version matches documentation
echo "Checking Svelte version..."
FRONTEND_SVELTE_VERSION=$(grep '"svelte":' "$FRONTEND_DIR/package.json" | grep -oP '[\d.]+' | head -1)
DOCS_SVELTE_VERSION=$(grep -oP "Svelte \K[\d.]+" "$DOCS_DIR/technology-stack.md" | head -1)

if [ "$FRONTEND_SVELTE_VERSION" != "$DOCS_SVELTE_VERSION" ]; then
    echo "  WARNING: package.json says Svelte $FRONTEND_SVELTE_VERSION but docs say $DOCS_SVELTE_VERSION"
    echo "  (Minor version differences may be acceptable)"
else
    echo "  OK: Svelte version $FRONTEND_SVELTE_VERSION"
fi

# 3. Check Vite version matches documentation
echo "Checking Vite version..."
FRONTEND_VITE_VERSION=$(grep '"vite":' "$FRONTEND_DIR/package.json" | grep -oP '[\d.]+' | head -1)
DOCS_VITE_VERSION=$(grep -oP "Vite \K[\d.]+" "$DOCS_DIR/technology-stack.md" | head -1)

if [ "$FRONTEND_VITE_VERSION" != "$DOCS_VITE_VERSION" ]; then
    echo "  WARNING: package.json says Vite $FRONTEND_VITE_VERSION but docs say $DOCS_VITE_VERSION"
else
    echo "  OK: Vite version $FRONTEND_VITE_VERSION"
fi

# 4. Validate internal markdown links
echo "Checking internal links..."
BROKEN_LINKS=0
for file in "$DOCS_DIR"/*.md; do
    # Extract relative links (not http/https)
    links=$(grep -oP '\]\((?!http)([^)]+)\)' "$file" 2>/dev/null | sed 's/\](\(.*\))/\1/' || true)
    for link in $links; do
        # Remove anchor
        link_file=$(echo "$link" | cut -d'#' -f1)
        if [ -n "$link_file" ]; then
            target="$DOCS_DIR/$link_file"
            if [ ! -f "$target" ] && [ ! -f "${target%.md}.md" ]; then
                echo "  BROKEN: $file -> $link"
                ((BROKEN_LINKS++))
            fi
        fi
    done
done

if [ $BROKEN_LINKS -eq 0 ]; then
    echo "  OK: All internal links valid"
else
    echo "  Found $BROKEN_LINKS broken links"
    ((ERRORS+=BROKEN_LINKS))
fi

# 5. Check for stale timestamps (older than 30 days)
echo "Checking documentation freshness..."
CURRENT_DATE=$(date +%s)
for file in "$DOCS_DIR"/*.md; do
    if [ -f "$file" ]; then
        FILE_DATE=$(stat -c %Y "$file" 2>/dev/null || stat -f %m "$file" 2>/dev/null)
        DAYS_OLD=$(( (CURRENT_DATE - FILE_DATE) / 86400 ))
        if [ $DAYS_OLD -gt 30 ]; then
            echo "  WARNING: $file is $DAYS_OLD days old"
        fi
    fi
done
echo "  Freshness check complete"

# 6. Count Svelte components and compare to docs
echo "Checking component count..."
ACTUAL_COMPONENTS=$(find "$FRONTEND_DIR/src/lib/components" -name "*.svelte" 2>/dev/null | wc -l)
echo "  Found $ACTUAL_COMPONENTS Svelte components in src/lib/components"

# Summary
echo ""
echo "=== Summary ==="
if [ $ERRORS -eq 0 ]; then
    echo "All checks passed!"
    exit 0
else
    echo "Found $ERRORS error(s)"
    exit 1
fi
