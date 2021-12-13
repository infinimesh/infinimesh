  FROM golang:1.11 AS builder
  ADD go.mod /go/src/github.com/slntopp/infinimesh/go.mod
  ADD . /go/src/github.com/slntopp/infinimesh
  WORKDIR /go/src/github.com/slntopp/infinimesh
  RUN cd /go/src/github.com/slntopp/infinimesh && CGO_ENABLED=0 go build ./cmd/avro-connector/

  FROM scratch
  WORKDIR /
  COPY --from=builder  /go/src/github.com/slntopp/infinimesh/avro-connector /avro-connector
  ENTRYPOINT ["/avro-connector"]
