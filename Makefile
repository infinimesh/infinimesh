
build-all:
	sh hack/build_images.sh

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
	
.PHONY: build-all mocks