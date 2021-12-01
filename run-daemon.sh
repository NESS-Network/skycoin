#!/usr/bin/env bash

# Runs Privateness in server daemon configuration

set -x

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
echo "Ness binary dir:" "$DIR"
pushd "$DIR" >/dev/null

COMMIT=$(git rev-parse HEAD)
BRANCH=$(git rev-parse --abbrev-ref HEAD)
GOLDFLAGS="${GOLDFLAGS} -X main.Commit=${COMMIT} -X main.Branch=${BRANCH}"

GORUNFLAGS=${GORUNFLAGS:-}
go run -ldflags "${GOLDFLAGS}" $GORUNFLAGS cmd/privateness/privateness.go \
    -gui-dir="${DIR}/src/gui/static/" \
    -max-default-peer-outgoing-connections=1 \
    -launch-browser=false \
    -enable-all-api-sets=true \
    -enable-gui=false \
    -log-level=info \
    -disable-pex=false \
    $@

popd >/dev/null
