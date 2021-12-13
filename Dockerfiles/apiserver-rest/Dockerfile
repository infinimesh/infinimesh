FROM golang:1.17-alpine AS builder

RUN apk add upx

WORKDIR /go/src/github.com/slntopp/infinimesh
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" ./cmd/apiserver-rest
RUN upx ./apiserver-rest

FROM alpine
WORKDIR /
COPY --from=builder  /go/src/github.com/slntopp/infinimesh/apiserver-rest /apiserver-rest

LABEL org.opencontainers.image.source https://github.com/slntopp/infinimesh

ENTRYPOINT ["/apiserver-rest"]
