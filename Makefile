build:
	GO111MODULE=on go build ./cmd/mqtt-bridge/
build_docker:
	docker build -t infinimesh/mqtt-bridge:latest -f Dockerfiles/mqtt-bridge/Dockerfile .
