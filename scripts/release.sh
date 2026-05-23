#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

if [ -n "$(git status --porcelain | grep -v '^?? ' || true)" ]; then
    echo "Working tree has uncommitted changes. Commit or stash first."
    git status --short
    exit 1
fi

CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
echo "==> Branch: $CURRENT_BRANCH"

if git rev-parse --abbrev-ref --symbolic-full-name @{u} >/dev/null 2>&1; then
    echo "==> Pulling latest"
    git pull --ff-only
fi

echo "==> Syncing OpenAPI spec from production"
make sync-spec

SPEC_CHANGED=0
if ! git diff --quiet gen/spec.json; then
    SPEC_CHANGED=1
    echo "    Spec changed:"
    git diff --stat gen/spec.json | sed 's/^/    /'
else
    echo "    Spec unchanged."
fi

echo "==> Regenerating and verifying"
make build >/dev/null
make test  >/dev/null

LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
VERSION=$(echo "$LATEST_TAG" | awk -F. '{
    if (NF<3) { print $0".0.1"; next }
    $NF=$NF+1; OFS="."; print
}')

if git rev-parse "$VERSION" >/dev/null 2>&1; then
    echo "Tag $VERSION already exists. Aborting."
    exit 1
fi

if [ $SPEC_CHANGED -eq 1 ]; then
    git add gen/spec.json
    git commit -m "chore: sync OpenAPI spec for $VERSION"
fi

git tag -a "$VERSION" -m "Release $VERSION"

if git rev-parse --abbrev-ref --symbolic-full-name @{u} >/dev/null 2>&1; then
    git push --follow-tags
else
    git push -u origin "$CURRENT_BRANCH" --follow-tags
fi

echo
echo "==> Released $VERSION (was $LATEST_TAG)"
echo "    GitHub Actions is building binaries and updating the Homebrew tap."
if command -v gh >/dev/null 2>&1; then
    echo "    Watch: gh run watch"
fi
