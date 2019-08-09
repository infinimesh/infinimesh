# Infinimesh Platform
Infinimesh is an opinionated multi-tenant hyperscale Platform to connect IoT devices securely. It features a unique Graph-based authorization system, allowing users & engineers to create arbitrary hierarchical ontologies, with the possibility to scope permissions down to single sub-devices to specific users (e.g. suppliers). It exposes simple to consume RESTful & gRPC APIs with both high-level (e.g. device shadow) and low-level (sending messages) concepts. Infinimesh Platform is open source and fully kubernetes compliant. No vendor lock-in - **run it yourself on Kubernetes in your own datacenter, under your control with maximum data privacy.**

## Project status
Infinimesh is currently under open source development. All development, is open source and completely transparent on GitHub. APIs are alpha and may change at any time. Many components are already available: MQTT Bridge, State management, Graph-based permission management, Device-to-Cloud and Cloud-to-Device messages. A kubernetes operator is also available, which is in an early stage. The simplest way to work with infinimesh is using a kubernetes based development environment: (https://github.com/infinimesh/infinimesh/tree/master/hack/microk8s).

A ui/dashboard is currently under development and will be available in mid Q2 2019. 
<br /> [Here](https://github.com/infinimesh/infinimesh/blob/master/roadmap.md) is a link to our feature roadmap.

## Build status
[![CircleCI](https://img.shields.io/circleci/project/github/infinimesh/infinimesh.svg)](https://circleci.com/gh/infinimesh/infinimesh/tree/master) [![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Finfinimesh%2Finfinimesh.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Finfinimesh%2Finfinimesh?ref=badge_shield)
[![GoReportCard](https://goreportcard.com/badge/github.com/infinimesh/infinimesh)](https://goreportcard.com/report/github.com/infinimesh/infinimesh)

| Docker Image  | Build status  |
| ------------- |---------------|
| Kubernetes Operator | [![Docker Repository on Quay](https://quay.io/repository/infinimesh/operator/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/operator) |
| App (Web UI) | [![Docker Repository on Quay](https://quay.io/repository/infinimesh/frontend/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/frontend) |
| API Server | [![Docker Repository on Quay](https://quay.io/repository/infinimesh/apiserver-rest/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/apiserver-rest) [![Docker Repository on Quay](https://quay.io/repository/infinimesh/apiserver/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/apiserver) |
| Node Server | [![Docker Repository on Quay](https://quay.io/repository/infinimesh/nodeserver/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/nodeserver) |
| Device Registry | [![Docker Repository on Quay](https://quay.io/repository/infinimesh/device-registry/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/device-registry) |
| Telemetry Router | [![Docker Repository on Quay](https://quay.io/repository/infinimesh/telemetry-router/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/telemetry-router) |
| MQTT-Bridge | [![Docker Repository on Quay](https://quay.io/repository/infinimesh/mqtt-bridge/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/mqtt-bridge) |
| Shadow | [![Docker Repository on Quay](https://quay.io/repository/infinimesh/shadow-delta-merger/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/shadow-delta-merger) [![Docker Repository on Quay](https://quay.io/repository/infinimesh/shadow-api/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/shadow-api) [![Docker Repository on Quay](https://quay.io/repository/infinimesh/shadow-persister/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/shadow-persister) |

## API Documentation
You can find swagger docs for the API server [here](https://infinimesh.github.io/infinimesh/swagger-ui/)

## Community
You can reach out to the community via [Slack](https://launchpass.com/infinimeshcommunity)

## Development
### Local development installation
We have built an automated local development setup based on microk8s:
```
bash <(curl -s https://raw.githubusercontent.com/infinimesh/infinimesh/master/hack/microk8s/infinimesh-setup.sh)
```
Source: https://github.com/infinimesh/infinimesh/tree/master/hack/microk8s

### Generating proto files
```
npm i -g merge-yaml-cli
npm i -g api-spec-converter
```

Start the local dev environment via `docker-compose up -d`.

Load test data into the database via `go run hack/import_dgraph_sample_data.go`

Login: `curl -X POST -d '{"username" : "joe", "password": "test123"}'  localhost:8081/account/token`

Get Objects: `curl -H 'Authorization: Bearer YOURTOKEN' localhost:8081/objects`

Login locally via CLI:
```
inf config set-context local --apiserver localhost:8080 --tls=false
inf login
```
Use the users joe / test123 or admin/admin123 for local development.

Register a device:
```
inf device create sample-device --cert-file hack/device_certs/sample_1.crt
```

Send sample message to the local instance:
```
mosquitto_pub --cafile hack/server.crt   --cert hack/device_certs/sample_1.crt --key hack/device_certs/sample_1.key -m '{"sensor" : {"temp" : 41}}' -t "devices/0x6ddd1/state/reported/delta" -h localhost  --tls-version tlsv1.2 -d -p 8089
```

Remember to replace 0x6ddd1 with the ID of your device. Also use the certificate and key of your device.

Send sample message via `mosquitto_pub` to the hosted SaaS instance:
```
mosquitto_pub --cafile /etc/ssl/certs/ca-certificates.crt   --cert hack/server.crt --key hack/server.key -m "blaaa" -t "shadows/testdeviceX" -h mqtt.api.infinimesh.io  --tls-version tlsv1.2 -d -p 8883
```

Access the frontend at http://localhost:8082

The cafile path may vary depending on your operating system.

## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Finfinimesh%2Finfinimesh.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Finfinimesh%2Finfinimesh?ref=badge_large)
