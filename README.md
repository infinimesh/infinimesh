# infinimesh IoT Platform

infinimesh is a opinionated multi-tenant hyperscale Internet of Things platform to connect IoT devices fast and securely with minimal TCO. It features a unique Graph-based authorization system, allowing users & engineers to create arbitrary hierarchical ontologies, with the possibility to scope permissions down to single sub-devices to specific users (e.g. suppliers). It exposes simple to consume RESTful & gRPC APIs with both high-level (e.g. device shadow) and low-level (sending messages) concepts. The infinimesh IoT platform is open source and fully kubernetes compliant. No vendor lock-in - **run it yourself on Kubernetes in your own datacenter, under your control with maximum data privacy.**

Our API's (REST / gRPC) are considered als beta and may change in future. infinimesh has already available:  
  
- **MQTT support for version 3 and 5**
- **State management (digital twin)**  
- **Graph-based permission management (multi-dimensional permissons at data layer)**  
- **TLS 1.2 / 1.3 support**  
- **Device-to-Cloud and Cloud-to-Device messages**  
- **Enhanced UI**  
- **k8s and Docker environments**

## Documentation  

Check out our [Wiki here](https://github.com/infinimesh/infinimesh/wiki).

## Build status

[![CI(Build Docker images)](https://github.com/infinimesh/infinimesh/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/infinimesh/infinimesh/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/infinimesh/infinimesh)](https://goreportcard.com/report/github.com/infinimesh/infinimesh)

## Community

You can reach out to the community via Discord.

![](http://invidget.switchblade.xyz/801798988163448832)

## Development

### Local development installation

We have built an automated local development setup based on Docker.

1. Add this entries to `/etc/hosts`:

  ```hosts
  127.0.0.1 api.infinimesh.local 
  127.0.0.1 console.infinimesh.local
  127.0.0.1 traefik.infinimesh.local
  127.0.0.1 rbmq.infinimesh.local
  127.0.0.1 db.infinimesh.local
  127.0.0.1 media.infinimesh.local
  127.0.0.1 mqtt.infinimesh.local
  ```

2. Close this repo
3. Run `docker compose up`

### Generating proto files

Clone [proto repo](https://github.com/infinimesh/proto)

Navigate to cloned repo directory and run:

```shell
docker run -it \
  -v $(pwd):/go/src/github.com/infinimesh/proto \
  ghcr.io/infinimesh/proto/buf:latest
```

Right now we keep protos generated only for Go. If you need one of the other languages, add according module to `buf.gen.yaml`.

PRs are as always welcome.

### Local Development

Start the local dev environment via `docker compose up -d`.

Load test data into the database via `go run hack/import_dgraph_sample_data.go`

Login locally via CLI:

```shell
inf login api.infinimesh.local infinimesh infinimesh --insecure
```

Access the Console at <http://console.infinimesh.local>

## License

Copyright 2018 - 2022, The Infinite Devices team

Licensed under the Apache License, Version 2.0 (the "Licenses"); you may not use
this file except in compliance with the License. You may obtain a copy of the License at

  [https://github.com/infinimesh/infinimesh/blob/master/LICENSE](https://github.com/infinimesh/infinimesh/blob/master/LICENSE)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

The authors of infinimesh are Infinite Devices GmbH, [birdayz](https://github.com/birdayz) and [2pk03](https://github.com/2pk03), all rights reserved.
