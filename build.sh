#!/bin/bash

set -ex

output=${1:-enapter}

BUILD_VERSION=$(git describe --tag 2> /dev/null)
BUILD_COMMIT=$(git rev-parse --short HEAD)
BUILD_DATE=$(date +'%Y-%m-%dT%H:%M:%S')

go build \
    -ldflags="-X 'main.version=${BUILD_VERSION}' -X 'main.commit=${BUILD_COMMIT}' -X 'main.date=${BUILD_DATE}'" \
    -o "${output}" ./cmd/enapter
