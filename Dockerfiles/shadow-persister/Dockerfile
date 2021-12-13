FROM golang:1.17-alpine AS builder

RUN apk add upx

WORKDIR /go/src/github.com/slntopp/infinimesh
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" ./cmd/shadow-persister
RUN upx ./shadow-persister

FROM scratch
WORKDIR /
COPY --from=builder /go/src/github.com/slntopp/infinimesh/shadow-persister /shadow-persister

LABEL org.opencontainers.image.source https://github.com/slntopp/infinimesh

ENTRYPOINT ["/shadow-persister"]
