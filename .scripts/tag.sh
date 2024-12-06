#!/bin/bash
set -euo pipefail

# Ensure commands work in non-interactive environments
export TERM=dumb

BOLD=$(tput bold)
NORMAL=$(tput sgr0)
echo "${BOLD}Tagging${NORMAL}"

# Get highest tag number, exit if no tags exist
VERSION=$(git describe --abbrev=0 --tags 2>/dev/null) || {
    echo "Error: No tags found in repository" >&2
    exit 1
}

# Validate semantic version format
if ! [[ $VERSION =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Error: Latest tag '$VERSION' doesn't match semantic versioning format (X.Y.Z)" >&2
    exit 1
}

# Split version into array safely
IFS='.' read -ra VERSION_BITS <<< "$VERSION"
# Validate version components
if [ ${#VERSION_BITS[@]} -ne 3 ]; then
    echo "Error: Version must have exactly three components (X.Y.Z)" >&2
    exit 1
fi

# Support different version increment types
INCREMENT_TYPE=${1:-patch}  # Default to patch increment if not specified
VNUM1=${VERSION_BITS[0]}
VNUM2=${VERSION_BITS[1]}
VNUM3=${VERSION_BITS[2]}

case $INCREMENT_TYPE in
    major)
        VNUM1=$((VNUM1+1))
        VNUM2=0
        VNUM3=0
        ;;
    minor)
        VNUM2=$((VNUM2+1))
        VNUM3=0
        ;;
    patch)
        VNUM3=$((VNUM3+1))
        ;;
    *)
        echo "Error: Invalid increment type. Use 'major', 'minor', or 'patch'" >&2
        exit 1
        ;;
esac
#create new tag
NEW_TAG="$VNUM1.$VNUM2.$VNUM3"

echo "Updating $VERSION to $NEW_TAG"

# Get current hash and see if it already has a tag
GIT_COMMIT=$(git rev-parse HEAD)
NEEDS_TAG=$(git describe --contains "$GIT_COMMIT" 2>/dev/null || true)

# Add --dry-run flag support
DRY_RUN=${DRY_RUN:-false}

if [ -n "$NEEDS_TAG" ]; then
    echo "Error: Current commit already has tag: $NEEDS_TAG" >&2
    exit 1
fi

# Show what will be done
echo "Ready to:"
echo "- Create tag: $NEW_TAG"
echo "- Push to remote"

if [ "$DRY_RUN" = "true" ]; then
    echo "Dry run complete"
    exit 0
fi

# Prompt for confirmation in interactive mode
if [ -t 0 ]; then
    read -p "Proceed? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Operation cancelled"
        exit 1
    fi
fi

# Create and push tag
if ! git tag "$NEW_TAG"; then
    echo "Error: Failed to create tag" >&2
    exit 1
fi

if ! git push --tags; then
    echo "Error: Failed to push tags" >&2
    git tag -d "$NEW_TAG"  # Cleanup on failure
    exit 1
fi