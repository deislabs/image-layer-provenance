#!/bin/bash

set -aueo pipefail

if [ -z "$1" ]; then
    echo "Builds and pushes all images in the examples directory to the registry."
    echo "Usage: $0 <registry-url>, e.g. $0 myregistry.registry.io"
    exit 1
fi

registry_url="$1"
namespace="$registry_url/image-manifest-layer-history"

example_image_directory_names=(
    "examples/postgres-base-multistage"
    "examples/postgres-base-singlestage"
    "examples/postgres-base-vanilla"
    "examples/scratch-base-multistage"
    "examples/scratch-base-singlestage"
    "examples/scratch-base-vanilla"
)

for example_image_directory_name in "${example_image_directory_names[@]}"; do
    image_ref="$namespace/$example_image_directory_name:latest"

    echo "[*] Building and pushing $image_ref"

    docker build -t "$image_ref" "$example_image_directory_name"
    docker push "$image_ref"
done
