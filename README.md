# infinimesh IoT Platform

infinimesh is a opinionated multi-tenant hyperscale Internet of Things platform to connect IoT devices fast and securely with minimal TCO. It features a unique Graph-based authorization system, allowing users & engineers to create arbitrary hierarchical ontologies, with the possibility to scope permissions down to single sub-devices to specific users (e.g. suppliers). It exposes simple to consume RESTful & gRPC APIs with both high-level (e.g. device shadow) and low-level (sending messages) concepts. The infinimesh IoT platform is open source and fully kubernetes compliant. No vendor lock-in - **run it yourself on Kubernetes in your own datacenter, under your control with maximum data privacy.**

Our API's (REST / gRPC / ConnectRPC) are considered as beta and may change in future. infinimesh has already available:  
  
- **MQTT support for version 3 and 5**
- **State management (digital twin)**  
- **Graph-based permission management (multi-dimensional permissons at data layer)**  
- **TLS 1.2 / 1.3 support**  
- **Device-to-Cloud and Cloud-to-Device messages**  
- **Enhanced UI**  
- **k8s and Docker environments**

## Documentation  

Check out our:

- [Wiki here](https://github.com/infinimesh/infinimesh/wiki).
- [Swagger UI](https://infinimesh.github.io/infinimesh/) **soon to be deprecated**

## Build status

[![CI(Build Docker images)](https://github.com/infinimesh/infinimesh/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/infinimesh/infinimesh/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/infinimesh/infinimesh)](https://goreportcard.com/report/github.com/infinimesh/infinimesh)

## Community

You can reach out to the community via Discord.

![](http://invidget.switchblade.xyz/801798988163448832)

## Client libraries

The recommended way to interact with infinimesh API is gRPC and ConnectRPC. Thus, we provide protobuf files at [`infinimesh/proto`](https://github.com/infinimesh/proto), using which you can generate a client library for your programming language of choice.

Additionally, we have pregenerated libraries:

### Golang

Can be added to your project via:

```shell
go get github.com/infinimesh/proto@latest
```

### JS

Can be obtained from npm via:

```shell
npm install infinimesh-proto
```

### Dart / FLutter

Has to be downloaded from [`infinimesh/proto`](https://github.com/infinimesh/proto), you can copy the whole [`build/dart`](https://github.com/infinimesh/proto/tree/master/build/dart) dir into your Flutter project

## CLI

### Usage

Start with `inf help` and `inf help login` ;)

### Homebrew

See [macOS](#macos).

### Snap

Just run

```shell
snap install inf
```

and see usage [usage](#usage)

### Linux

#### `.deb` (Debian, Ubuntu, etc.)

1. Go to [CLI Releases](https://github.com/infinimesh/inf/releases)
2. Get `.deb` package for your CPU arch (`arm64` or `x86_64`)
3. `dpkg -i path/to/.deb`

If you're using some other arch, let us know, we'll add it to the build. Meanwhile - try [building from source](#build-from-source)

Then see usage [usage](#usage)

#### `.rpm` (RedHat, CentOS, Fedora, etc.)

1. Go to [CLI Releases](https://github.com/infinimesh/inf/releases)
2. Get `.rpm` package for your CPU arch (`arm64` or `x86_64`)
3. `yum localinstall path/to/.rpm` or `dnf install path/to/.rpm`

If you're using some other arch, let us know, we'll add it to the build. Meanwhile - try [building from source](#build-from-source)

Then see usage [usage](#usage)

#### AUR (Arch Linux, Manjaro, etc.)

If you have `yaourt` or `yay` package must be found automatically by label `inf-bin`

Otherwise,

1. `git clone https://aur.archlinux.org/packages/inf-bin`
2. `cd inf-bin`
3. `makepkg -i`

Then see usage [usage](#usage)

#### Others

If you're using other package manager or have none, you can download prebuilt binary in `.tar.gz` archive for `arm64` or `x86_64`, unpack it and put `inf` binary to `/usr/bin` or your `$PATH/bin`.

If you're using some other arch, let us know, we'll add it to the build. Meanwhile - try [building from source](#build-from-source)

Then see usage [usage](#usage)

### macOS

If you're using [**Homebrew**](https://brew.sh):

```shell
brew tap infinimesh/inf
brew install inf
```

You're good to go!

Then see usage [usage](#usage)

If you don't have [**Homebrew**](https://brew.sh), consider using it ;), otherwise you can get prebuilt binary from [CLI Releases page](https://github.com/infinimesh/inf/releases) as an `.tar.gz` archive.

```shell
# if you have wget then
wget https://github/infinimesh/inf/releases/#version/inf-version-darwin-arch.tar.gz
# if you don't, just download it
tar -xvzf #inf-version-darwin-arch.tar.gz
# move binary to /usr/local/bin or alike
mv #inf-version-darwin-arch/inf /usr/local/bin
```

You're good to go!

Then see usage [usage](#usage)

### Windows

1. Go to [CLI Releases](https://github.com/infinimesh/inf/releases)
2. Get prebuilt binary from [CLI Releases page](https://github.com/infinimesh/inf/releases) as an `.zip` archive.
3. Unpack it
4. Put it somewhere in `$PATH`

Then see usage [usage](#usage)

### Build From Source

See [CLI repo](https://github.com/infinimesh/inf) for source and instructions.

## Development

### Local development installation

#### Production like

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

2. Clone this repo via `git clone https://github.com/infinimesh/infinimesh.git` 
3. cd into the fresh cloned repo
4. Copy `.env.example` into `.env`, edit if needed
5. Run `docker compose up`

Swagger API: https://infinimesh.github.io/infinimesh/

#### Debugging

This repo contains VS Code `launch.json` commands you can use to spin up `infinimesh` locally with debugger attached.

Requirements:

- Docker with Docker Compose module
- Traefik
- Bash Debug Extention
- Go Extention
- **Hosts entries from the above must be present**

How to launch:

1. Go to Run and Debug
2. Click `Lauch Debug Environment` | [!NOTE] You'll be asked the root password to launch `traefik`, don't miss it
3. Click `Spin the infinimesh up` to launch all services
4. Enjoy

This will run all services locally (not in Docker and with Debuggers attached). Since about everything needs a port or two, here are the reserved ports and hostnames:

|       Service             |    Ports   |
|---------------------------|------------|
| api.infinimesh.local      | 80/http    |
| db.infinimesh.local       | 80/http    |
| console.infinimesh.local  | 80/http    |
| traefik.infinimesh.local  | 80/http    |
| rbmq.infinimesh.local     | 80/http    |
| MQTT Bridge (Basic Auth)  | 1883/http  |
| RabbitMQ API              | 5672/tcp   |
| Redis                     | 6379/tcp   |
| Redis Timeseries          | 6380/tcp   |
| Console Dev server        | 5173/http  |
| api.infinimesh.local      | 8000/grpc  |
| Node (repo)               | 8001       |
| REST API Gateway          | 8002       |
| Shadow                    | 8003/grpc  |
| _TimeSeries API_          | 8004       |
| ArangoDB                  | 8529/http  |
| MQTT Bridge               | 8883/https |
| RabbitMQ UI               | 15672/http |

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

### Mocks

infinimesh is using [mockery](https://vektra.github.io/mockery). To generate mocks it is recommended to use Docker image:

```shell
docker pull vektra/mockery
docker run -v "$PWD":/src -w /src vektra/mockery --all
```

Otherwise just get the `mockery` binary and run it at the root of the project.

### Local Development

Start the local dev environment via `docker compose up -d`.

Login locally via CLI:

```shell
inf login api.infinimesh.local infinimesh infinimesh --insecure
```

Access the Console at <http://console.infinimesh.local>

## License

Copyright 2018 - 2024 the infinimesh committers, 2pk03, birdayz, slntopp and The Infinite AI Audio GmbH team

Licensed under the Apache License, Version 2.0 (the "Licenses"); you may not use
this file except in compliance with the License. You may obtain a copy of the License at

  [https://github.com/infinimesh/infinimesh/blob/master/LICENSE](https://github.com/infinimesh/infinimesh/blob/master/LICENSE)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

The authors of infinimesh are The Infinite AI Audio GmbH, [birdayz](https://github.com/birdayz), [2pk03](https://github.com/2pk03) and [slntopp](https://github.com/slntopp), all rights reserved.
