
build-all:
	sh hack/build_images.sh

mocks:
	docker run -v "$PWD":/src -w /src vektra/mockery --all