/*
Copyright © 2022 Johnson Shi <Johnson.Shi@microsoft.com>
*/
package history

const (
	// FROMPrimaryBaseImageLayer is the layer creation type for
	// primary base image layers (layers created by `FROM <primary-base-image>`).
	FROMPrimaryBaseImageLayer LayerCreationType = "FROM-PrimaryBaseImageLayer"

	// COPYFromMultistageBuildStageLayer is the layer creation type for
	// COPY from multistage build stage layers (layers created by `COPY --from=<multistage-build-stage> <src> <dst>`).
	COPYFromMultistageBuildStageLayer LayerCreationType = "COPY-FromMultistageBuildStageLayer"

	// COPYCommandLayer is the layer creation type for
	// for non-multistage COPY commands (layers created by `COPY <src> <dst>`).
	COPYCommandLayer LayerCreationType = "COPY-CommandLayer"

	// ADDCommandLayer is the layer creation type for
	// for layers created by `ADD`.
	ADDCommandLayer LayerCreationType = "ADD-CommandLayer"

	// RUNCommandLayer is the layer creation type for
	// for layers created by `RUN`.
	RUNCommandLayer LayerCreationType = "RUN-CommandLayer"

	// UNKNOWNDockerfileCommandLayer is the layer creation type for
	// for layers created by unknown Dockerfile commands.
	// This should not happen as this would imply that commands other than
	// {FROM, COPY, ADD, RUN} created a layer, which is impossible.
	UNKNOWNDockerfileCommandLayer LayerCreationType = "UNKNOWN-DockerfileCommandLayer"
)
