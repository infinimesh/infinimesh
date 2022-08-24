FROM golang:1.19-alpine AS builder

RUN apk add upx

WORKDIR /go/src/github.com/infinimesh/infinimesh
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -buildvcs=false ./cmd/mqtt-bridge
RUN upx ./mqtt-bridge

FROM scratch
WORKDIR /
COPY --from=builder /go/src/github.com/infinimesh/infinimesh/mqtt-bridge /mqtt-bridge

LABEL org.opencontainers.image.source https://github.com/infinimesh/infinimesh

ENTRYPOINT ["/mqtt-bridge"]
