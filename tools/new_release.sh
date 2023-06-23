#!/bin/bash
set -e

version=$1

if [ -z "$version" ]; then
    echo "Usage: ./tools/new_release.sh 0.2.1"
    exit 1
fi

prev_release=$(git tag|tail -1)
date=$(date '+%B %-d, %Y')
changelog=$(git log --pretty=format:%s "$prev_release"..HEAD | sed 's/^/\* /g')
echo -e "## $version ($date)\n$changelog\n" | cat - CHANGELOG.md > CHANGELOG.temp && mv CHANGELOG.temp CHANGELOG.md

echo -e "✅ All release files generated.\nPlease review it and commit with desired commit message:\n\n  git commit -m \"release: v$version\""
