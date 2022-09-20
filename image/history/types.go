/*
Copyright © 2022 Johnson Shi <Johnson.Shi@microsoft.com>
*/
package history

import (
	"github.com/asottile/dockerfile"
	"github.com/docker/docker/api/types/image"
	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// ImageHistory describes the overall history of an image.
// It contains:
// - the image's non-detailed layer history,
//   - This is the non-detailed layer history from Docker Engine API's ImageHistory operation (`docker image history`).
// - the image's OCI image manifest,
// - the image's Dockerfile.
type ImageHistory struct {
	// ImageLayerHistory MUST be sorted from top layers (most recent layers) to bottom layers (base image layers).
	ImageLayerHistory []image.HistoryResponseItem `json:"ImageLayerHistory"`
	// OCI image manifest. See https://github.com/opencontainers/image-spec/blob/main/manifest.md.
	ImageManifest ocispec.Manifest `json:"ImageManifest"`
	// DockerfileCommands MUST be sorted based on the original order of commands in the Dockerfile.
	// E.g. if the Dockerfile contains the following commands:
	// 		FROM ubuntu:22.04
	// 		RUN apt-get update
	// 		RUN apt-get install -y vim
	// then dockerfileCommands should be:
	// 		[]dockerfile.Command{
	// 			dockerfile.Command{Original: "FROM ubuntu:22.04", Cmd: "FROM", Value: []string{"ubuntu:22.04"}},
	// 			dockerfile.Command{Original: "RUN apt-get update", Cmd: "RUN", Value: []string{"apt-get", "update"}},
	// 			dockerfile.Command{Original: "RUN apt-get install -y vim", Cmd: "RUN", Value: []string{"apt-get", "install", "-y", "vim"}},
	// 		}
	DockerfileCommands []dockerfile.Command `json:"DockerfileCommands"`
}

// LayerCreationType describes the type of layer creation.
type LayerCreationType string

// LayerCreationParameters describes the exact Dockerfile parameters used to create the image manifest layer.
type LayerCreationParameters struct {
	DockerfileLayerCreationType LayerCreationType    `json:"DockerfileLayerCreationType"`
	DockerfileCommands          []dockerfile.Command `json:"DockerfileCommands"`
	BaseImageRef                string               `json:"BaseImage"`
}

// ImageManifestLayerDockerfileCommandsHistory describes the exact Dockerfile history of an image manifest layer.
type ImageManifestLayerDockerfileCommandsHistory struct {
	LayerDescriptor         ocispec.Descriptor      `json:"LayerDescriptor"`
	LayerCreationParameters LayerCreationParameters `json:"LayerCreationParameters"`
	AttributedEntity        map[string]string       `json:"AttributedEntity"`
}

// SimplifiedImageManifestLayerDockerfileCommandsHistory describes the exact Dockerfile history
// of an image manifest layer (simplified format).
type SimplifiedImageManifestLayerDockerfileCommandsHistory struct {
	LayerDigest                 digest.Digest     `json:"LayerDigest"`
	DockerfileLayerCreationType LayerCreationType `json:"DockerfileLayerCreationType"`
	DockerfileCommands          []string          `json:"DockerfileCommands"`
	BaseImageRef                string            `json:"BaseImage"`
	AttributedEntity            map[string]string `json:"AttributedEntity"`
}
