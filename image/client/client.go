package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
)

type ImageClient struct {
	Username           string
	Password           string
	DockerEngineClient *client.Client
	ImageRef           registry.Reference
	OrasRepoClient     *remote.Repository
}

// NewImageClient returns a new ImageClient.
func NewImageClient(username, password, imageRefStr string) (ImageClient, error) {
	imageClient := ImageClient{
		Username: username,
		Password: password,
	}

	// Create a client to the Docker Engine (Docker daemon) referenced by local environment variables.
	dockerEngineClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return imageClient, err
	}
	imageClient.DockerEngineClient = dockerEngineClient

	// Parse the image reference string to an ORAS artifact reference type.
	imageRef, err := registry.ParseReference(imageRefStr)
	if err != nil {
		return imageClient, err
	}
	imageClient.ImageRef = imageRef

	// Create an ORAS client to the image reference's repository.
	orasRepoClient, err := remote.NewRepository(imageClient.ImageRef.String())
	if err != nil {
		return imageClient, err
	}
	orasRepoClient.Client = &auth.Client{
		Credential: func(ctx context.Context, reg string) (auth.Credential, error) {
			return auth.Credential{
				Username: imageClient.Username,
				Password: imageClient.Password,
			}, nil
		},
	}
	imageClient.OrasRepoClient = orasRepoClient

	return imageClient, nil
}

// EnsureImageIsPulled ensures that the image is pulled
// and present in the Docker Engine (Docker daemon).
func (imageClient *ImageClient) EnsureImageIsPulled(ctx context.Context) error {
	// Create an encoded auth config for authenticating with the image reference's registry.
	// See https://docs.docker.com/engine/api/sdk/examples/#pull-an-image-with-authentication.
	authConfig := types.AuthConfig{
		Username: imageClient.Username,
		Password: imageClient.Password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return err
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	out, err := imageClient.DockerEngineClient.ImagePull(
		ctx,
		imageClient.ImageRef.String(),
		types.ImagePullOptions{RegistryAuth: authStr},
	)
	if err != nil {
		return err
	}
	defer out.Close()

	// Fully read and consume to ensure the image is pulled
	// to the Docker Engine (Docker daemon).
	_, err = ioutil.ReadAll(out)
	if err != nil {
		// If reading the output reader fails,
		// then most likely pulling the image failed.
		return err
	}

	return nil
}

// GetImageHistory returns the image layer history.
// The image layer history is sorted
// from top layers (most recent layers) to bottom layers (base image layers).
// See https://docs.docker.com/engine/reference/commandline/image_history/.
func (imageClient *ImageClient) GetImageLayerHistory(ctx context.Context) ([]image.HistoryResponseItem, error) {
	// The image must be present in the Docker Engine (Docker daemon)
	// before its history can be obtained.
	err := imageClient.EnsureImageIsPulled(ctx)
	if err != nil {
		return nil, err
	}

	// The image history can only be obtained once the image has been pulled
	// and is present in the Docker Engine (Docker daemon).
	images, err := imageClient.DockerEngineClient.ImageList(ctx, types.ImageListOptions{
		All: true,
		Filters: filters.NewArgs(
			filters.Arg("reference", imageClient.ImageRef.String()),
		),
	})
	if err != nil {
		return nil, err
	}
	if len(images) == 0 {
		return nil, fmt.Errorf("no image found for reference: %s", imageClient.ImageRef.String())
	}
	if len(images) > 1 {
		return nil, fmt.Errorf("multiple images found for reference: %s", imageClient.ImageRef.String())
	}

	imageID := images[0].ID
	imageLayerHistory, err := imageClient.DockerEngineClient.ImageHistory(ctx, imageID)
	if err != nil {
		return nil, err
	}

	return imageLayerHistory, nil
}

// GetImageManifest returns the OCI image manifest.
// See https://github.com/opencontainers/image-spec/blob/main/manifest.md.
func (imageClient *ImageClient) GetImageManifest(ctx context.Context) (ocispec.Manifest, error) {
	manifest := ocispec.Manifest{}

	descriptor, rc, err := imageClient.OrasRepoClient.FetchReference(ctx, imageClient.ImageRef.String())
	if err != nil {
		return manifest, err
	}
	defer rc.Close()

	pulledBytes, err := content.FetchAll(ctx, imageClient.OrasRepoClient, descriptor)
	if err != nil {
		return manifest, err
	}
	if descriptor.Size != int64(len(pulledBytes)) || descriptor.Digest != digest.FromBytes(pulledBytes) {
		return manifest, fmt.Errorf("pulled bytes do not match descriptor")
	}

	err = json.Unmarshal(pulledBytes, &manifest)
	if err != nil {
		return manifest, err
	}

	return manifest, nil
}
