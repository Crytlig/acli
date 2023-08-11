.PHONY: vendor
vendor:
	GO111MODULE=on go mod vendor;

.PHONY: build
build:
	GO111MODULE=on go build -v -o .

.PHONY: build-all
build-all: vendor build-windows-amd64 build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64

# TODO: Fix go build. Seems like the main.Version is not overwritten
.PHONY: build-windows-amd64
build-windows-amd64:
	GOOS=windows GOARCH=amd64 go build -ldflags '-X "github.com/crytlig/acli/version.Version=${ACLI_VERSION}"' -v -o ./bin/acli-${ACLI_VERSION}-windows-amd64.exe

.PHONY: build-linux-amd64
build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -ldflags '-X "github.com/crytlig/acli/version.Version=${ACLI_VERSION}"' -v -o ./bin/acli-${ACLI_VERSION}-linux-amd64

.PHONY: build-linux-arm64
build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -ldflags '-X "github.com/crytlig/acli/version.Version=${ACLI_VERSION}"' -v -o ./bin/acli-${ACLI_VERSION}-linux-arm64

.PHONY: build-darwin-amd64
build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -ldflags '-X "github.com/crytlig/acli/version.Version=${ACLI_VERSION}"' -v -o ./bin/acli-${ACLI_VERSION}-darwin-amd64

.PHONY: build-darwin-arm64
build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -ldflags '-X "github.com/crytlig/acli/version.Version=${ACLI_VERSION}"' -v -o ./bin/acli-${ACLI_VERSION}-darwin-arm64

