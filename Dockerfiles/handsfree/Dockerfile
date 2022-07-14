FROM golang:1.18-alpine AS builder

RUN apk add upx

WORKDIR /go/src/github.com/infinimesh/infinimesh
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -buildvcs=false ./cmd/handsfree
RUN upx ./handsfree

FROM scratch
WORKDIR /
COPY --from=builder /go/src/github.com/infinimesh/infinimesh/handsfree /handsfree

LABEL org.opencontainers.image.source https://github.com/infinimesh/infinimesh

ENTRYPOINT ["/handsfree"]
