VERSION ?= "latest"

build-all:
	sh hack/build_images.sh

build-console:
	export INFINIMESH_VERSION_TAG=$(git describe --tags --abbrev=0)
	docker build . -f "Dockerfiles/console/Dockerfile" -t "ghcr.io/infinimesh/infinimesh/console:${VERSION}"

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