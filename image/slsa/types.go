package slsa

import (
	"time"

	"github.com/johnsonshi/image-manifest-layer-history/image/history"
)

type ImageManifestLayerSlsaProvenance struct {
	LayerHistory                 history.ImageManifestLayerDockerfileCommandsHistory `json:"LayerHistory"`
	BuilderID                    string                                              `json:"BuilderID"`
	BuildType                    string                                              `json:"BuildType"`
	BuildInvocationID            string                                              `json:"BuildInvocationID"`
	BuildStartedOn               *time.Time                                          `json:"BuildStartedOn"`
	BuildFinishedOn              *time.Time                                          `json:"BuildFinishedOn"`
	RepoURIContainingImageSource string                                              `json:"RepoURIContainingImageSource"`
	RepoGitCommit                string                                              `json:"RepoGitCommit"`
	RepoPathToImageSource        string                                              `json:"RepoPathToImageSource"`
}
