/*
Copyright © 2022 Johnson Shi <Johnson.Shi@microsoft.com>
*/
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"time"

	"github.com/asottile/dockerfile"
	intoto "github.com/in-toto/in-toto-golang/in_toto"
	"github.com/spf13/cobra"

	"github.com/johnsonshi/image-manifest-layer-history/pkg/image/client"
	"github.com/johnsonshi/image-manifest-layer-history/pkg/image/history"
	imagehistoryslsa "github.com/johnsonshi/image-manifest-layer-history/pkg/image/history/slsa"
)

type generateCmdOpts struct {
	stdin                    io.Reader
	stdout                   io.Writer
	stderr                   io.Writer
	username                 string
	password                 string
	imageRef                 string
	dockerfilePath           string
	outputFilePath           string
	simplifiedJsonOutput     bool
	slsaProvenanceJsonOutput bool
	attributionAnnotations   []string
}

func newGenerateCmd(stdin io.Reader, stdout io.Writer, stderr io.Writer, args []string) *cobra.Command {
	opts := &generateCmdOpts{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}

	cobraCmd := &cobra.Command{
		Use:     "generate",
		Short:   "Generate the history (including the exact Dockerfile commands) for each OCI Image Manifest Layer of a container image.",
		Example: `generate -u username -p password -i imageRef -d dockerfilePath -o outputFilePath [-s] [-a "key: value"]`,
		RunE: func(_ *cobra.Command, args []string) error {
			return opts.run()
		},
	}

	f := cobraCmd.Flags()

	f.StringVarP(&opts.username, "username", "u", "", "username to use for authentication with the registry")
	cobraCmd.MarkFlagRequired("username")

	// TODO add support for --password-stdin (reading password from stdin) for more secure password input.
	f.StringVarP(&opts.password, "password", "p", "", "password to use for authentication with the registry")
	cobraCmd.MarkFlagRequired("password")

	cobraCmd.MarkFlagsRequiredTogether("username", "password")

	f.StringVarP(&opts.imageRef, "image-ref", "i", "", "full image reference including registry, repository, and tag/digest, e.g. myregistry.azurecr.io/library/ubuntu:22.04 or myregistry.azurecr.io/library/ubuntu@sha256:123456")
	cobraCmd.MarkFlagRequired("image-ref")

	f.StringVarP(&opts.dockerfilePath, "dockerfile", "d", "", "path to the Dockerfile")
	cobraCmd.MarkFlagRequired("dockerfile")

	f.StringVarP(&opts.outputFilePath, "output-file", "o", "", "optional path to an output file")

	f.BoolVar(&opts.simplifiedJsonOutput, "simplified-json", false, "optional flag to output in simplified JSON format")

	f.BoolVar(&opts.slsaProvenanceJsonOutput, "slsa-provenance-json", false, "optional flag to output in SLSA Provenance JSON format")

	cobraCmd.MarkFlagsMutuallyExclusive("simplified-json", "slsa-provenance-json")

	f.StringArrayVarP(&opts.attributionAnnotations, "attribution-annotation", "a", []string{}, "optional flag to add 'string-key:string-value' attributions to the manifest layer history (only added for layers not 'FROM <primary-base-image>')")

	return cobraCmd
}

func (opts *generateCmdOpts) run() error {
	ctx := context.Background()

	annotationsMap, err := getAnnotationsMap(opts.attributionAnnotations)
	if err != nil {
		return err
	}

	imageClient, err := client.NewImageClient(opts.username, opts.password, opts.imageRef)
	if err != nil {
		return err
	}

	// imageLayerHistory is sorted from top layers (most recent layers) to bottom layers (base image layers).
	imageLayerHistory, err := imageClient.GetImageLayerHistory(ctx)
	if err != nil {
		return err
	}

	imageManifest, err := imageClient.GetImageManifest(ctx)
	if err != nil {
		return err
	}

	// dockerfileCommands is sorted based on the original order of commands in the Dockerfile.
	// E.g. if the Dockerfile contains the following commands:
	// 		FROM ubuntu:22.04
	// 		RUN apt-get update
	// 		RUN apt-get install -y vim
	// then dockerfileCommands will be:
	// 		[]dockerfile.Command{
	// 			dockerfile.Command{Original: "FROM ubuntu:22.04", Cmd: "FROM", Value: []string{"ubuntu:22.04"}},
	// 			dockerfile.Command{Original: "RUN apt-get update", Cmd: "RUN", Value: []string{"apt-get", "update"}},
	// 			dockerfile.Command{Original: "RUN apt-get install -y vim", Cmd: "run", Value: []string{"apt-get", "install", "-y", "vim"}},
	// 		}
	dockerfileCommands, err := dockerfile.ParseFile(opts.dockerfilePath)
	if err != nil {
		return err
	}

	h := history.ImageHistory{
		ImageLayerHistory:  imageLayerHistory,
		ImageManifest:      imageManifest,
		DockerfileCommands: dockerfileCommands,
	}

	manifestLayerHistory, err := h.GetImageManifestLayerDockerfileCommandsHistory(annotationsMap)
	if err != nil {
		return err
	}

	return opts.writeManifestLayerHistory(manifestLayerHistory)

}

func (opts *generateCmdOpts) writeManifestLayerHistory(manifestLayerHistory []history.ImageManifestLayerDockerfileCommandsHistory) error {
	output, err := json.MarshalIndent(manifestLayerHistory, "", "  ")
	if err != nil {
		return err
	}

	if opts.simplifiedJsonOutput {
		simplified := history.GetSimplifiedImageManifestLayerDockerfileCommandsHistory(manifestLayerHistory)
		output, err = json.MarshalIndent(simplified, "", "  ")
		if err != nil {
			return err
		}
	} else if opts.slsaProvenanceJsonOutput {
		var imageSlsaProvenanceStatements []*intoto.ProvenanceStatement
		timeNow := time.Now()
		for _, layerHistory := range manifestLayerHistory {
			layerSlsaProvenance := imagehistoryslsa.ImageManifestLayerSlsaProvenance{
				LayerHistory:                layerHistory,
				BuilderID:                   "URI indicating the builder identity. E.g. pipeline-name",
				BuildType:                   "URI indicating what type of build was performed. E.g. build-type-dockerfile-build",
				BuildInvocationID:           "Globally Unique Build Invocation ID. Definition: Identifies this particular build invocation, which can be useful for finding associated logs or other ad-hoc analysis. The exact meaning and format is defined by builder.id; by default it is treated as opaque and case-sensitive. The value SHOULD be globally unique.",
				BuildStartedOn:              &timeNow,
				BuildFinishedOn:             &timeNow,
				RepoURIContainingDockerfile: "URI to Git repo of Dockerfile. Describes where the config file that kicked off the build came from. URI indicating the identity of the source of the config. E.g. https://www.github.com/example/reponame/blob/master/Dockerfile",
				RepoGitCommit:               "Git commit SHA that kicked off the build.",
				RepoPathToDockerfile:        "Path to Dockerfile in the repo. Definition: String identifying the entry point into the build. This is often a path to a configuration file and/or a target label within that file. The syntax and meaning are defined by buildType. For example, if the buildType were “make”, then this would reference the directory in which to run make as well as which target to use.",
			}
			layerSlsaProvenanceStatement, err := layerSlsaProvenance.GetImageManifestLayerSlsaProvenance()
			if err != nil {
				return err
			}
			imageSlsaProvenanceStatements = append(imageSlsaProvenanceStatements, layerSlsaProvenanceStatement)
		}
		output, err = json.MarshalIndent(imageSlsaProvenanceStatements, "", "  ")
		if err != nil {
			return err
		}
	}

	output = append(output, '\n')

	if opts.outputFilePath != "" {
		return writeToFile(opts.outputFilePath, output)
	}

	_, err = opts.stdout.Write(output)
	return err
}

// getAnnotationsMap returns a map of annotations from a slice of annotation strings.
// strings in the slice should conform to the following format: "key: value".
func getAnnotationsMap(annotationSlice []string) (map[string]string, error) {
	re := regexp.MustCompile(`:\s*`)
	annotationsMap := make(map[string]string)
	for _, rawAnnotation := range annotationSlice {
		annotation := re.Split(rawAnnotation, 2)
		if len(annotation) != 2 {
			return nil, fmt.Errorf("invalid annotation: %s", rawAnnotation)
		}
		annotationsMap[annotation[0]] = annotation[1]
	}
	return annotationsMap, nil
}

func writeToFile(path string, data []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}
