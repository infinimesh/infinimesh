include .env

VERSION ?= "latest"
GIT_VERSION ?= $(shell git describe --tags --abbrev=0)
GIT_COMMIT ?= $(shell git rev-parse --short HEAD)

build-all:
	sh hack/build_images.sh

build-console-dry:
	@echo "Building console image:"
	@echo "docker build . -f \"Dockerfiles/console/Dockerfile\" --build-arg=\"INFINIMESH_VERSION_TAG=${GIT_VERSION}\" --build-arg=\"INFINIMESH_COMMIT_HASH=${GIT_COMMIT}\" -t \"ghcr.io/infinimesh/infinimesh/console:${VERSION}\""

build-console:
	docker build . -f "Dockerfiles/console/Dockerfile" --build-arg="INFINIMESH_VERSION_TAG=${GIT_VERSION}" --build-arg="INFINIMESH_COMMIT_HASH=${GIT_COMMIT}" -t "ghcr.io/infinimesh/infinimesh/console:${VERSION}"

build-web:
	docker build . -f "Dockerfiles/web/Dockerfile" -t "ghcr.io/infinimesh/infinimesh/web:${VERSION}"

build-repo:
	docker build . -f "Dockerfiles/repo/Dockerfile" -t "ghcr.io/infinimesh/infinimesh/repo:${VERSION}"

build-mqtt:
	docker build . -f "Dockerfiles/mqtt-bridge/Dockerfile" -t "ghcr.io/infinimesh/infinimesh/mqtt-bridge:${VERSION}"

build-shadow:
	docker build . -f "Dockerfiles/shadow/Dockerfile" -t "ghcr.io/infinimesh/infinimesh/shadow:${VERSION}"

mocks:
	docker run -v "$PWD":/src -w /src vektra/mockery --all

GOBIN ?= $$(go env GOPATH)/bin

.PHONY: install-go-test-coverage
install-go-test-coverage:
	go install github.com/vladopajic/go-test-coverage/v2@latest

.PHONY: check-coverage
check-coverage: install-go-test-coverage
	go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
	${GOBIN}/go-test-coverage --config=./.testcoverage.yml

check-coverage-html: check-coverage
	rm cover.html 2> /dev/null
	go tool cover -html=cover.out -o=cover.html

.PHONY: build-all build-console mocks

vscode:
	@docker compose -f vscode.docker-compose.yaml up -d
	@traefik --configfile vscode.traefik.yml

vscode-logs:
	@docker compose -f vscode.docker-compose.yaml logs -f