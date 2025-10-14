#!/bin/bash

set -eo pipefail

if [[ $# -ne 1 ]]; then
  echo "Usage: ${0} [patch|minor|major]"
  exit 1
fi

BUMP_TYPE="${1}"

case "${BUMP_TYPE}" in
  patch|minor|major)
    ;;
  *)
    echo "Invalid bump type: ${BUMP_TYPE}. Use patch, minor, or major."
    exit 1
    ;;
esac

LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "0.0.0")

IFS='.' read -r MAJOR MINOR PATCH <<< "${LATEST_TAG#v}"

case "${BUMP_TYPE}" in
  "patch")
    PATCH=$((PATCH + 1))
    ;;
  "minor")
    MINOR=$((MINOR + 1))
    PATCH=0
    ;;
  "major")
    MAJOR=$((MAJOR + 1))
    MINOR=0
    PATCH=0
    ;;
esac

NEW_VERSION="${MAJOR}.${MINOR}.${PATCH}"

git tag "v${NEW_VERSION}"
git push origin "v${NEW_VERSION}"
