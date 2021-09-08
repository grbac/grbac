#!/usr/bin/env sh

set -o errexit
set -o nounset
set -o pipefail

sleep 10

grbac init --dgraph-endpoint=dgraph:9080
grbac run --dgraph-endpoint=dgraph:9080

exit 0