GOLANGCI_LINT_VERSION = 2.5.0
ACTIONLINT_VERSION = 1.7.7

IMG_NAME ?= 'ghcr.io/epicstep/gdatum'
IMG_VERSION ?= '$(shell whoami)'

generate:
	go generate ./...

build:
	GOOS=linux go build -o gdatum cmd/app/main.go
	docker build --build-arg IMG_TAG=debug-nonroot -f Containerfile -t $(IMG_NAME):$(IMG_VERSION) .
	rm -rf gdatum

test:
	go test --timeout 10m -race ./...

coverage:
	go test -race -v -coverpkg=./... -coverprofile=profile.out ./...
	go tool cover -func profile.out

lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v$(GOLANGCI_LINT_VERSION) run

actionlint:
	go run github.com/rhysd/actionlint/cmd/actionlint@v$(ACTIONLINT_VERSION)
