FROM golang:1.17-alpine AS builder

RUN apk add upx

WORKDIR /go/src/github.com/slntopp/infinimesh
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" ./cmd/nodeserver
RUN upx ./nodeserver

FROM scratch
WORKDIR /
COPY --from=builder /go/src/github.com/slntopp/infinimesh/nodeserver /nodeserver

LABEL org.opencontainers.image.source https://github.com/slntopp/infinimesh

ENTRYPOINT ["/nodeserver"]
