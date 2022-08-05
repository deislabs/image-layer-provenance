/*
Copyright © 2022 Johnson Shi <Johnson.Shi@microsoft.com>
*/
package slsa

import (
	"time"

	"github.com/johnsonshi/image-manifest-layer-history/image/history"
)

type ImageManifestLayerSlsaProvenance struct {
	LayerHistory                history.ImageManifestLayerDockerfileCommandsHistory `json:"LayerHistory"`
	BuilderID                   string                                              `json:"BuilderID"`
	BuildType                   string                                              `json:"BuildType"`
	BuildInvocationID           string                                              `json:"BuildInvocationID"`
	BuildStartedOn              *time.Time                                          `json:"BuildStartedOn"`
	BuildFinishedOn             *time.Time                                          `json:"BuildFinishedOn"`
	RepoURIContainingDockerfile string                                              `json:"RepoURIContainingDockerfile"`
	RepoGitCommit               string                                              `json:"RepoGitCommit"`
	RepoPathToDockerfile        string                                              `json:"RepoPathToDockerfile"`
}
