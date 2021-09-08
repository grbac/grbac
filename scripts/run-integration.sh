#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

export INTEGRATION_TEST_DGRAPH_ENDPOINT=127.0.0.1:9080

# Launch the dgraph docker container and open its ports.
echo "integration: starting the dgraph docker container..."
container_id=$(docker run --detach --rm -p 9080:9080 dgraph/standalone:v21.03.0)

# Wait for the container to be up and running.
echo "integration: waiting (10s) for the container to be ready..."
sleep 10s

# Run the integration tests and store the return code of the 'go test' command.
go test -cover -tags=integration ./... && return_code=$? || return_code=$?    

# Stop the dgraph docker container.
echo "integration: stopping the container..."
docker stop $container_id

exit $return_code