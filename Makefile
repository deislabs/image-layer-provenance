.PHONY: build-cli
build-cli:
	go build -v -o ./bin/image-layer-dockerfile-history ./cmd/cli
	chmod +x ./bin/image-layer-dockerfile-history
