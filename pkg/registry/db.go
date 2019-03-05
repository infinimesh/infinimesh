package registry

import (
	"github.com/infinimesh/infinimesh/pkg/node/dgraph"
)

type Device struct {
	dgraph.Object
	Tags    []string `json:"tags,omitempty"`
	Enabled bool     `json:"enabled"`

	Certificates []*X509Cert `json:"certificates,omitempty"`
}

type X509Cert struct {
	dgraph.Node
	PemData              string `json:"pem_data,omitempty"`
	Algorithm            string `json:"algorithm,omitempty"`
	Fingerprint          []byte `json:"fingerprint,omitempty"`
	FingerprintAlgorithm string `json:"fingerprint.algorithm,omitempty"`
}
