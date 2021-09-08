#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

WORKDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

echo "updating go modules..."

GOPROXY=direct go get -u github.com/animeapis/api-go-client@master
GOPROXY=direct go get -u github.com/animeapis/go-genproto@master

echo "updating git submodules..."

git submodule foreach git pull origin master

echo "regenerating gapics..."

source "${WORKDIR}/gapic.sh"

exit 0