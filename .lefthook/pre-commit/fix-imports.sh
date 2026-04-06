#!/bin/bash
# Auto-format staged Go files before lint check

STAGED_FILES=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$')

if [ -z "$STAGED_FILES" ]; then
  exit 0
fi

echo "$STAGED_FILES" | xargs go fmt

# Re-add formatted files
echo "$STAGED_FILES" | xargs git add

exit 0
