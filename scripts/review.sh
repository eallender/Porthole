#!/bin/bash

# Professional code review script for pre-commit validation
# Uses go-code-reviewer agent for comprehensive analysis

set -e

# Color definitions for clean output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
print_header() {
    echo -e "${BLUE}============================================${NC}"
    echo -e "${BLUE} Pre-Commit Code Review${NC}"
    echo -e "${BLUE}============================================${NC}"
}

print_section() {
    echo -e "\n${BLUE}[$1]${NC}"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}!${NC} $1"
}

# Main execution
print_header

# Check if there are staged changes
if git diff --cached --quiet; then
    print_error "No staged changes found"
    echo "Please stage files before running review:"
    echo "  git add <files>"
    exit 1
fi

# Show staged files
print_section "STAGED FILES"
git diff --cached --name-only | sed 's/^/  /'

# Run Go-specific checks if Go files are staged
if git diff --cached --name-only | grep -q '\.go$'; then
    print_section "GO TOOLING CHECKS"
    
    echo -n "  go fmt... "
    if go fmt ./... >/dev/null 2>&1; then
        echo -e "${GREEN}PASS${NC}"
    else
        echo -e "${RED}FAIL${NC}"
        print_error "go fmt failed. Please fix formatting issues."
        exit 1
    fi
    
    echo -n "  go vet... "
    if go vet ./... >/dev/null 2>&1; then
        echo -e "${GREEN}PASS${NC}"
    else
        echo -e "${RED}FAIL${NC}"
        print_error "go vet failed. Please fix issues."
        exit 1
    fi
    
    echo -n "  go mod tidy... "
    if go mod tidy >/dev/null 2>&1; then
        echo -e "${GREEN}PASS${NC}"
    else
        echo -e "${RED}FAIL${NC}"
        print_error "go mod tidy failed."
        exit 1
    fi
    
    # Check if go mod tidy changed anything
    changed_files=""
    if ! git diff --quiet go.mod 2>/dev/null; then
        changed_files="go.mod"
    fi
    if [ -f go.sum ] && ! git diff --quiet go.sum 2>/dev/null; then
        changed_files="$changed_files go.sum"
    fi
    
    if [ -n "$changed_files" ]; then
        print_warning "go mod tidy made changes to:$changed_files"
        echo "Please stage the updated files:"
        echo "  git add$changed_files"
        exit 1
    fi
fi

# AI Code Review
print_section "AI CODE REVIEW"
echo "Running production-quality analysis..."
echo ""

# Create a temporary prompt file for the staged changes
TEMP_PROMPT=$(mktemp)
cat > "$TEMP_PROMPT" << 'EOF'
Please review the staged changes in this commit for production readiness as a Go code reviewer. Focus on:
- Critical issues that would prevent deployment
- TODO comments and code quality issues 
- Important improvements needed
- Security, performance, and Go best practices

Provide a concise, well-formatted review suitable for terminal output.

Here are the staged changes:
EOF

# Append the actual diff to the prompt
git diff --cached >> "$TEMP_PROMPT"

# Call claude with the prompt file
/usr/local/bin/claude --print < "$TEMP_PROMPT"

# Clean up
rm "$TEMP_PROMPT"

# Completion
echo ""
echo -e "${BLUE}============================================${NC}"
print_success "Review completed successfully"
echo ""
echo "Next steps:"
echo "  • Address any feedback above"
echo "  • Re-run review: ./scripts/review.sh"
echo "  • Commit changes: git commit"
echo "  • Or unstage: git reset"
echo -e "${BLUE}============================================${NC}"