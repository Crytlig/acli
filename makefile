.PHONY: vendor
vendor:
	GO111MODULE=on go mod vendor;

.PHONY: build
build:
	GO111MODULE=on go build -v -o .

.PHONY: build-all
build-all: go-generate vendor build-windows-amd64 build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64

.PHONY: build-windows-amd64
build-windows-amd64:
	GOOS=windows GOARCH=amd64 go build -ldflags "-X github.com/Crytlig/acli/main.Version=${ACLI_VERSION}" -v -o ./bin/draft-windows-amd64.exe

.PHONY: build-linux-amd64
build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/Crytlig/acli/main.Version=${ACLI_VERSION}" -v -o ./bin/draft-linux-amd64

.PHONY: build-linux-arm64
build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -ldflags "-X github.com/Crytlig/acli/main.Version=${ACLI_VERSION}" -v -o ./bin/draft-linux-arm64

.PHONY: build-darwin-amd64
build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X github.com/Crytlig/acli/main.Version=${ACLI_VERSION}" -v -o ./bin/draft-darwin-amd64

.PHONY: build-darwin-arm64
build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X github.com/Crytlig/acli/main.Version=${ACLI_VERSION}" -v -o ./bin/draft-darwin-arm64
