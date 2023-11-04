#!/bin/sh

set -eu

export VERSION=${VERSION:-0.0.0}
export GOFLAGS="'-ldflags=-w -s \"-X=github.com/jcdickinson/ollama/version.Version=$VERSION\" \"-X=github.com/jcdickinson/ollama/server.mode=release\"'"

docker buildx build \
    --load \
    --platform=linux/arm64,linux/amd64 \
    --build-arg=VERSION \
    --build-arg=GOFLAGS \
    -f Dockerfile \
    -t ollama \
    .
