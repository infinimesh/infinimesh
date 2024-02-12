FROM node:16 AS console-builder

ARG INFINIMESH_VERSION_TAG
ENV INFINIMESH_VERSION_TAG=$INFINIMESH_VERSION_TAG

ARG INFINIMESH_COMMIT_HASH
ENV INFINIMESH_COMMIT_HASH=$INFINIMESH_COMMIT_HASH

ADD console /console
WORKDIR /console

RUN npm install -g pnpm
RUN pnpm i && pnpm build

# Build server
FROM golang:1.20-alpine AS builder

RUN apk add upx

WORKDIR /go/src/github.com/infinimesh/infinimesh

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd/console ./cmd/console
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -buildvcs=false ./cmd/console
RUN upx ./console

# Switch to the scratch image
FROM scratch

COPY --from=console-builder /console/dist/ /app
COPY --from=builder /go/src/github.com/infinimesh/infinimesh/console /console

EXPOSE 80

LABEL org.opencontainers.image.source https://github.com/infinimesh/infinimesh

# Run thttpd
ENTRYPOINT ["/console"]