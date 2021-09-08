#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

API_NAME="grbac"
API_VERSION="v1alpha1"

# TODO: Everything should be moved to Bazel for protobuf compilation.

# Generate CLI via GAPIC.
protoc \
  --experimental_allow_proto3_optional \
  --proto_path="schema/api-common-protos" \
  --proto_path="schema/animeapis" \
  --go_cli_out="cmd" \
  --go_cli_opt="root=grbac" \
  --go_cli_opt="gapic=github.com/animeapis/api-go-client/${API_NAME}/${API_VERSION}" \
  --go_cli_opt="fmt=true" \
    "schema/animeapis/animeshon/${API_NAME}/${API_VERSION}/${API_NAME}.proto"

exit 0