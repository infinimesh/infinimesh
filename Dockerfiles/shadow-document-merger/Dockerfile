FROM golang:1.11 AS builder
ADD go.mod /go/src/github.com/infinimesh/infinimesh/go.mod
ADD . /go/src/github.com/infinimesh/infinimesh
WORKDIR /go/src/github.com/infinimesh/infinimesh
RUN cd /go/src/github.com/infinimesh/infinimesh && CGO_ENABLED=0 GO111MODULE=on go build ./cmd/shadow-document-merger/

FROM scratch
WORKDIR /
COPY --from=builder  /go/src/github.com/infinimesh/infinimesh/shadow-document-merger /shadow-document-merger
ENTRYPOINT ["/shadow-document-merger"]
