#!/usr/bin/env bash

cd "$(dirname "${BASH_SOURCE[0]}")/../.."
set -exo pipefail

# Download sourcegraph cli
curl -L https://sourcegraph.com/.api/src-cli/src_linux_amd64 -o /usr/local/bin/src
chmod +x /usr/local/bin/src

# Download lsif-go
curl -L  https://github.com/sourcegraph/lsif-go/releases/download/v1.2.0/src_linux_amd64 -o /usr/local/bin/lsif-go
chmod +x /usr/local/bin/lsif-go


asdf install

IMAGE=sourcegraph/server:insiders ./dev/run-server-image.sh

pushd internal/cmd/precise-code-intel-tester

./scripts/download.sh

go build && ./precise-code-intel-tester upload

go build && ./precise-code-intel-tester query

