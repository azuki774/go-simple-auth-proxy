SHELL=/bin/bash
container_name=go-simple-auth-proxy

.PHONY: bin build test
bin:
	go build -a -tags "netgo" -installsuffix netgo  -ldflags="-s -w -extldflags \"-static\" \
	-X main.version=$(git describe --tag --abbrev=0) \
	-X main.revision=$(git rev-list -1 HEAD) \
	-X main.build=$(git describe --tags)" \
	-o ./bin/ ./...

build:
	docker build -t $(container_name):latest -f build/Dockerfile .

test:
	go vet ./...
	go test -v ./...
