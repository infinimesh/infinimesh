# Infinimesh Platform
Infinimesh Platform is an opinionated Platform to connect IoT devices securely. It exposes simple to consume RESTful & gRPC APIs with both high-level (e.g. device shadow) and low-level (sending messages) concepts. Infinimesh Platform is open source and fully cloud native. No vendor lock-in - run it yourself on Kubernetes or use our SaaS offering (TBA).

![asciicinema](asciicast.gif)

## Build status
[![CircleCI](https://img.shields.io/circleci/project/github/infinimesh/infinimesh.svg)](https://circleci.com/gh/infinimesh/infinimesh/tree/master) [![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Finfinimesh%2Finfinimesh.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Finfinimesh%2Finfinimesh?ref=badge_shield)
[![GoReportCard](https://goreportcard.com/badge/github.com/infinimesh/infinimesh)](https://goreportcard.com/report/github.com/infinimesh/infinimesh) 

| Docker Image  | Build status  |
| ------------- |---------------|
| API Server | [![Docker Repository on Quay](https://quay.io/repository/infinimesh/apiserver-rest/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/apiserver-rest) [![Docker Repository on Quay](https://quay.io/repository/infinimesh/apiserver/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/apiserver) |
| Controlplane-UI | [![Docker Repository on Quay](https://quay.io/repository/infinimesh/controlplane-ui/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/controlplane-ui) |
| Node Server | [![Docker Repository on Quay](https://quay.io/repository/infinimesh/nodeserver/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/nodeserver) |
| Device Registry | [![Docker Repository on Quay](https://quay.io/repository/infinimesh/device-registry/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/device-registry) |
| Telemetry Router | [![Docker Repository on Quay](https://quay.io/repository/infinimesh/telemetry-router/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/telemetry-router) |
| MQTT-Bridge | [![Docker Repository on Quay](https://quay.io/repository/infinimesh/mqtt-bridge/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/mqtt-bridge) |
| Shadow | [![Docker Repository on Quay](https://quay.io/repository/infinimesh/shadow-delta-merger/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/shadow-delta-merger) [![Docker Repository on Quay](https://quay.io/repository/infinimesh/shadow-api/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/shadow-api) [![Docker Repository on Quay](https://quay.io/repository/infinimesh/shadow-persister/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/shadow-persister) |

## API Documentation
You can find swagger docs for the API server [here](https://infinimesh.github.io/infinimesh/swagger-ui/)

## Development
### Generating proto files
```
npm i -g merge-yaml-cli 
npm i -g api-spec-converter
```
Load test data into the database via `go run hack/import_dgraph_sample_data.go`

Login: `curl -X POST -d '{"username" : "joe", "password": "test123"}'  localhost:8081/token`

Get Objects: `curl -H 'Authorization: Bearer YOURTOKEN' localhost:8081/objects`

Send sample message via `mosquitto_pub`:
```
mosquitto_pub --cafile /etc/ssl/certs/ca-certificates.crt   --cert hack/server.crt --key hack/server.key -m "blaaa" -t "shadows/testdeviceX" -h mqtt.api.infinimesh.io  --tls-version tlsv1.2 -d -p 8883
```

The cafile path may vary depending on your operating system.

## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Finfinimesh%2Finfinimesh.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Finfinimesh%2Finfinimesh?ref=badge_large)
