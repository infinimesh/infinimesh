# infinimesh IoT Platform
infinimesh is a opinionated multi-tenant hyperscale Internet of Things platform to connect IoT devices fast and securely with minimal TCO. It features a unique Graph-based authorization system, allowing users & engineers to create arbitrary hierarchical ontologies, with the possibility to scope permissions down to single sub-devices to specific users (e.g. suppliers). It exposes simple to consume RESTful & gRPC APIs with both high-level (e.g. device shadow) and low-level (sending messages) concepts. The infinimesh IoT platform is open source and fully kubernetes compliant. No vendor lock-in - **run it yourself on Kubernetes in your own datacenter, under your control with maximum data privacy.**

Our API's (REST / gRPC) are considered als beta and may change in future. infinimesh has already available:  
  
**MQTT support for version 3 and 5**   
**State management (digital twin)**  
**Graph-based permission management (multi-dimensional permissons at data layer)**  
**TLS 1.2 / 1.3 support**  
**Device-to-Cloud and Cloud-to-Device messages**  
**Enhanced UI**  
**k8s and docker environments**

A kubernetes operator is also available, which is in an early stage. The simplest way to work with infinimesh is using a kubernetes based development environment: (https://github.com/slntopp/infinimesh/tree/master/hack/microk8s). 

# Documentation  
Our [documentation](https://infinimesh.github.io/infinimesh/docs/#/) is getting better and better. Please file PR if you find mistakes or just want to add something. We review on daily basis.

## Build status
[![CircleCI](https://img.shields.io/circleci/project/github/infinimesh/infinimesh.svg)](https://circleci.com/gh/infinimesh/infinimesh/tree/master) 
[![Go Report Card](https://goreportcard.com/badge/github.com/slntopp/infinimesh)](https://goreportcard.com/report/github.com/slntopp/infinimesh)

## API Documentation
You can find swagger docs for the API server [here](https://infinimesh.github.io/infinimesh/swagger-ui/)

## Community
You can reach out to the community via [Slack](https://launchpass.com/infinimeshcommunity) or join us in our CNCF channel [#infinimesh](https://cloud-native.slack.com/archives/C01EP6QRJTD).

## Development
### Local development installation
We have built an automated local development setup based on microk8s.
For Ubuntu please use:
```
bash <(curl -s https://raw.githubusercontent.com/infinimesh/infinimesh/master/hack/microk8s/infinimesh-setup-ubuntu.sh)
```
For OSX please use:
```
bash <(curl -s https://raw.githubusercontent.com/infinimesh/infinimesh/master/hack/microk8s/infinimesh-setup-osx.sh)
```
Source: https://github.com/slntopp/infinimesh/tree/master/hack/microk8s

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

Copyright 2018 - 2021, The infinimesh team

Licensed under the Apache License, Version 2.0 (the "Licenses"); you may not use 
this file except in compliance with the License. You may obtain a copy of the License at

       https://github.com/slntopp/infinimesh/blob/master/LICENSE

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
