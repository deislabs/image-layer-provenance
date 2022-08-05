#!/bin/bash

set -aueo pipefail

if [ -z "$1" ] || [ -z "$2" ] || [ -z "$3" ]; then
    echo "Generates the OCI image manifest layer histories for all images referenced in the examples directory."
    echo "Usage: $0 <registry-url> <registry-username> <registry-password>"
    echo "  e.g. $0 myregistry.registry.io myregistryusername myregistrypassword"
    exit 1
fi

registry_url="$1"
namespace="$registry_url/image-manifest-layer-history"

username="$2"
password="$3"

example_image_directory_names=(
    "examples/postgres-base-multistage"
    "examples/postgres-base-singlestage"
    "examples/postgres-base-vanilla"
    "examples/scratch-base-multistage"
    "examples/scratch-base-singlestage"
    "examples/scratch-base-vanilla"
)

make build-cli

responsible_entity_id="ID of the responsible entity that introduced non-base image layers. i.e. Entity ID of the Dockerfile Author. This denotes responsible entity for the new layers added on top of base image layers."

for example_image_directory_name in "${example_image_directory_names[@]}"; do
    image_ref="$namespace/$example_image_directory_name:latest"
    dockerfile="$example_image_directory_name/Dockerfile"

    # Save a copy of the simple image history (i.e. the image history returned from `docker image history`).
    docker image history --no-trunc --human=false "$image_ref" > "$example_image_directory_name/simple-image-history.txt"

    # Save a copy of the OCI image manifest (i.e. the OCI image manifest returned from `docker manifest inspect`).
    docker manifest inspect "$image_ref" > "$example_image_directory_name/oci-image-manifest.json"

    echo -e "[*] Generating OCI image manifest layer history\n\tfor image: '$image_ref'\n\tbased on Dockerfile: '$dockerfile'\n"

    ./bin/image-layer-dockerfile-history \
        generate \
        --username "$username" \
        --password "$password" \
        --image-ref "$image_ref" \
        --dockerfile "$dockerfile" \
        --output-file "$example_image_directory_name/oci-image-manifest-layer-history.json" \
        --attribution-annotation "responsible_entity_id: $responsible_entity_id"

    ./bin/image-layer-dockerfile-history \
        generate \
        --username "$username" \
        --password "$password" \
        --image-ref "$image_ref" \
        --dockerfile "$dockerfile" \
        --output-file "$example_image_directory_name/oci-image-manifest-layer-history-simplified.json" \
        --simplified-json=true \
        --attribution-annotation "responsible_entity_id: $responsible_entity_id"

    ./bin/image-layer-dockerfile-history \
        generate \
        --username "$username" \
        --password "$password" \
        --image-ref "$image_ref" \
        --dockerfile "$dockerfile" \
        --output-file "$example_image_directory_name/oci-image-manifest-layer-history-slsa.json" \
        --slsa-provenance-json=true \
        --attribution-annotation "responsible_entity_id: $responsible_entity_id"
done
