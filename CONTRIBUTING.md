# Contributing to the Project

Welcome! We are very happy to accept community contributions to the project, whether those are [Pull Requests](#pull-requests), [Suggestions](#suggestions) or [Bug Reports](#bug-reports). Please note that by participating in this project, you agree to abide by the [Code of Conduct](./CODE_OF_CONDUCT.md), as well as the terms of the [CLA](#cla).

## Getting Started

* Follow the installation and demo steps to run the [Proof of Concept CLI Tool](./README.md#proof-of-concept-cli-tool).

## Developing

### Components

The project is composed of the following main components:

* **CLI** (`./cmd/cli`): the CLI executable.
* **Go SDK** (`./image`): Go SDK to generate container image layer history and provenance documents.
* **Scripts** (`./scripts`): scripts to build, generate, and push provenance documents for container images.

### Tests

* We are happy to accept contributions for additional tests in the CLI ([`./cmd/cli`](./cmd/cli)), the SDK ([`./image`](./image/)), and the scripts ([`./scripts`](./scripts/)).

### Running the CLI

* Once built, run the CLI from the bin directory `./bin/image-layer-dockerfile-history` for a list of the available commands.
* For any command the `--help` argument can be passed for more information and a list of possible arguments.

## Pull Requests

If you'd like to start contributing, you can search for GitHub Issues.

## Suggestions

* Please first search GitHub Issues before opening an issue to check whether your suggestion has already been suggested. If it has, feel free to add your own comments to the existing issue.
* Ensure you have included a "What?" - what your feature entails, being as specific as possible, and giving mocked-up syntax examples where possible.
* Ensure you have included a "Why?" - what the benefit of including this feature will be.

## Bug Reports

* Please first search GitHub Issues before opening an issue to see if it has already been reported.
* Try to be as specific as possible, including the version of the CLI or SDK used to reproduce the issue, and any example arguments needed to reproduce it.

## CLA

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.opensource.microsoft.com.

When you submit a pull request, a CLA bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., status check, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.
