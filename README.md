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
4. Run `docker compose up`

Swagger API: https://infinimesh.github.io/infinimesh/

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

Login locally via CLI:

```shell
inf login api.infinimesh.local infinimesh infinimesh --insecure
```

Access the Console at <http://console.infinimesh.local>

## License

Copyright 2018 - 2023, The Infinite Devices team

Licensed under the Apache License, Version 2.0 (the "Licenses"); you may not use
this file except in compliance with the License. You may obtain a copy of the License at

  [https://github.com/infinimesh/infinimesh/blob/master/LICENSE](https://github.com/infinimesh/infinimesh/blob/master/LICENSE)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

The authors of infinimesh are Infinite Devices GmbH, [birdayz](https://github.com/birdayz) and [2pk03](https://github.com/2pk03), all rights reserved.
