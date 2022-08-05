# Image Manifest Layer History

Command-line tool that shows the _**exact**_ Dockerfile commands for each [OCI Image Manifest](https://github.com/opencontainers/image-spec/blob/main/manifest.md) Layer of a container image.

## Quick Start

### Install

To install, run the following commands.

```bash
curl -LO https://github.com/johnsonshi/image-manifest-layer-history/releases/download/v0.0.1/image-layer-dockerfile-history
chmod +x image-layer-dockerfile-history
sudo mv image-layer-dockerfile-history /usr/local/bin
```

### Generate

Generate the history (including the exact Dockerfile commands) for each OCI Image Manifest Layer of a container image.

#### Generate – Usage

```bash
image-layer-dockerfile-history \
  generate \
  --username "$username" \
  --password "$password" \
  --image-ref "$image_ref" \
  --dockerfile "$dockerfile" \
  --output-file "$output_file" \
  --attribution-annotation "git_commit_id: $git_commit_id" \
  --attribution-annotation "git_commit_date: $git_commit_date" \
  --attribution-annotation "git_commit_name: $git_commit_name" \
  --attribution-annotation "git_commit_email: $git_commit_email" \
  --attribution-annotation "git_remote_origin_url: $git_remote_origin_url"
```

See [`./scripts/generate-history-all-examples.sh`](./scripts/generate-history-all-examples.sh) for more examples.
