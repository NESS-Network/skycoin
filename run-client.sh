#!/usr/bin/env bash

# Runs Privateness in desktop client configuration
export COIN=privateness
export DATA_DIR=/root/.privateness
export RPC_ADDR=http://127.0.0.1:6460
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
    -max-default-peer-outgoing-connections=7 \
    -launch-browser=true \
    -enable-all-api-sets=true \
    -enable-gui=true \
    -log-level=info \
    $@

popd >/dev/null
