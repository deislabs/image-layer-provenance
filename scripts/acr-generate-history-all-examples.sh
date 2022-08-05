#!/bin/bash

set -aueo pipefail

if [ -z "$1" ]; then
    echo "Generates the OCI image manifest layer histories for all images referenced in the examples directory."
    echo "Usage: $0 <ACR-Registry-Name>, e.g. $0 myregistry (not the URL, just the name)"
    exit 1
fi

acr_registry_name="$1"

acr_access_token_output=$(az acr login --name "$acr_registry_name" --expose-token)
acr_access_token_username="00000000-0000-0000-0000-000000000000"
acr_access_token=$(echo "$acr_access_token_output" | jq --raw-output ".accessToken")
acr_login_server=$(echo "$acr_access_token_output" | jq --raw-output ".loginServer")

echo "$acr_access_token" | docker login "$acr_login_server" -u "$acr_access_token_username" --password-stdin

./scripts/generate-history-all-examples.sh "$acr_login_server" "$acr_access_token_username" "$acr_access_token"
