package registry

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/wrappers"

	"github.com/infinimesh/infinimesh/pkg/node"
	"github.com/infinimesh/infinimesh/pkg/node/dgraph"
	"github.com/infinimesh/infinimesh/pkg/registry/registrypb"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"

	"encoding/base64"
)

type Server struct {
	dgo *dgo.Dgraph

	repo node.Repo
}

func NewServer(dg *dgo.Dgraph) *Server {
	return &Server{
		dgo: dg,
		repo: &dgraph.DGraphRepo{
			Dg: dg,
		},
	}
}

func (s *Server) getFingerprint(pemCert []byte, certType string) (fingerprint []byte, err error) {
	pemBlock, _ := pem.Decode(pemCert)
	if pemBlock == nil {
		return nil, errors.New("Could not decode PEM")
	}
	cert, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return sha256Sum(cert.Raw), nil
}

func sha256Sum(c []byte) []byte {
	s := sha256.New()
	_, err := s.Write(c)
	if err != nil {
		panic(err)
	}
	return s.Sum(nil)
}

func (s *Server) Create(ctx context.Context, request *registrypb.CreateRequest) (*registrypb.CreateResponse, error) {
	txn := s.dgo.NewTxn()
	defer txn.Discard(ctx) // nolint
	if exists := dgraph.NameExists(ctx, txn, request.Device.Name, request.Namespace, ""); exists {
		return nil, status.Error(codes.FailedPrecondition, "Name exists already")
	} // TODO allow setting parent

	ns, err := s.repo.GetNamespace(ctx, request.Namespace)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, "Invalid namespace")
	}

	if request.Device.Certificate == nil {
		return nil, status.Error(codes.FailedPrecondition, "No certificate provided")
	}

	fp, err := s.getFingerprint([]byte(request.Device.Certificate.PemData), request.Device.Certificate.Algorithm)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, "Invalid Certificate")
	}

	var enabled bool
	if request.Device.Enabled != nil {
		enabled = request.Device.Enabled.GetValue()
	}

	d := &Device{
		Object: dgraph.Object{
			Node: dgraph.Node{
				UID:  "_:new",
				Type: "object",
			},
			Name: request.Device.Name,
			Kind: node.KindDevice,
		},
		Enabled: enabled,
		Certificates: []*X509Cert{
			&X509Cert{
				PemData:              request.Device.Certificate.PemData,
				Algorithm:            request.Device.Certificate.Algorithm,
				Fingerprint:          fp,
				FingerprintAlgorithm: "sha256",
			},
		},
	}

	js, err := json.Marshal(d)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to create device")
	}
	mutRes, err := txn.Mutate(ctx, &api.Mutation{
		SetJson: js,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to create object: %v", err))
	}

	newUID := mutRes.GetUids()["new"]
	fmt.Println(newUID)

	nsMut := &api.NQuad{
		Subject:   ns.GetId(),
		Predicate: "owns",
		ObjectId:  newUID,
	}

	_, err = txn.Mutate(ctx, &api.Mutation{
		Set: []*api.NQuad{
			nsMut,
		},
		CommitNow: true,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to create object: %v", err))
	}

	return &registrypb.CreateResponse{
		Fingerprint: fp,
	}, nil
}

func (s *Server) Update(ctx context.Context, request *registrypb.UpdateRequest) (response *registrypb.UpdateResponse, err error) {
	d := &Device{
		Object: dgraph.Object{
			Node: dgraph.Node{
				UID: request.Device.Id,
			},
		},
	}
	for _, field := range request.FieldMask.GetPaths() {
		switch field {
		case "Enabled":
			d.Enabled = request.Device.Enabled.GetValue()
		case "Tags":
			d.Tags = request.Device.Tags
			//TODO
			// case "Certificate.Algorithm":
			// 	update["certificate_fingerprint_algorithm"] = request.Device.Certificate.Algorithm
			// case "Certificate.PemData":
			// 	update["certificate"] = request.Device.Certificate.PemData
		}

	}

	txn := s.dgo.NewTxn()
	js, err := json.Marshal(&d)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to patch device: %v", err))
	}

	_, err = txn.Mutate(ctx, &api.Mutation{
		SetJson:   js,
		CommitNow: true,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to patch device: %v", err))
	}

	// TODO
	// Updating cert currently not implemented

	// if _, ok := update["certificate"]; ok {
	// 	// recalc fingerprint
	// 	fp, err := s.getFingerprint([]byte(request.Device.Certificate.PemData), request.Device.Certificate.Algorithm)
	// 	if err != nil {
	// 		return nil, status.Error(codes.FailedPrecondition, "Invalid Certificate")
	// 	}

	// 	update["certificate_fingerprint"] = fp
	// 	update["certificate_fingerprint_algorithm"] = "sha256"
	// }

	// var device Device
	// if err := s.db.First(&device, "id = ?", request.Device.GetId()).Error; err != nil {
	// 	return nil, err
	// }

	// if err := s.db.Model(&device).Updates(update).Error; err != nil {
	// 	return nil, err
	// }

	return &registrypb.UpdateResponse{}, nil
}

func (s *Server) GetByFingerprint(ctx context.Context, request *registrypb.GetByFingerprintRequest) (*registrypb.GetByFingerprintResponse, error) {
	txn := s.dgo.NewReadOnlyTxn()

	const q = `query devices($fingerprint: string){
  devices(func: eq(fingerprint, "0r1H1PStHl9AY1WQlQVBxNqmq2FjkhlvBN9D9hDbOks=")) @normalize {
    ~certificates {
      uid : uid
      name : name
      ~owns {
        namespace: name
      }
    }
  }
}
  `

	vars := map[string]string{
		"$fingerprint": base64.StdEncoding.EncodeToString(request.Fingerprint),
	}

	resp, err := txn.QueryWithVars(ctx, q, vars)
	if err != nil {
		return nil, err
	}

	var res struct {
		Devices []struct {
			Uid       string `json:"uid"`
			Name      string `json:"name"`
			Namespace string `json:"namespace"`
		} `json:"devices"`
	}

	err = json.Unmarshal(resp.Json, &res)
	if err != nil {
		return nil, err
	}

	// FIXME dedupe, we could have multiple entries possibly?
	var devices []*registrypb.DeviceForFingerprint
	for _, device := range res.Devices {
		devices = append(devices, &registrypb.DeviceForFingerprint{
			Id:        device.Uid,
			Name:      device.Name,
			Namespace: device.Namespace,
		})
	}

	return &registrypb.GetByFingerprintResponse{
		Devices: devices,
	}, nil
}

func (s *Server) Get(ctx context.Context, request *registrypb.GetRequest) (response *registrypb.GetResponse, err error) {
	txn := s.dgo.NewReadOnlyTxn()

	const q = `query devices($name: string, $namespace: string){
                     devices(func: eq(name, $name)) @cascade {
                       uid
                       name
                       ~owns {
                       } @filter(eq(name, $namespace))
                       certificates {
                         pem_data
                         algorithm
                         fingerprint
                         fingerprint.algorithm
                       }
                     }
                   }`

	vars := map[string]string{
		"$name":      request.Id, // TODO rename id to name OR to device_id
		"$namespace": request.Namespace,
	}

	resp, err := txn.QueryWithVars(ctx, q, vars)
	if err != nil {
		return nil, err
	}

	var res struct {
		Devices []*Device `json:"devices"`
	}

	err = json.Unmarshal(resp.Json, &res)
	if err != nil {
		return nil, err
	}

	if len(res.Devices) == 0 {
		return &registrypb.GetResponse{}, status.Error(codes.NotFound, "Device not found")
	}

	return &registrypb.GetResponse{
		Device: toProto(res.Devices[0]),
	}, nil
}

func (s *Server) List(ctx context.Context, request *registrypb.ListDevicesRequest) (response *registrypb.ListResponse, err error) {
	txn := s.dgo.NewReadOnlyTxn()

	const q = `query list($account: string, $namespace: string){
                     var(func: eq(name,$namespace)) @filter(eq(type, "namespace")) {
                       owns {
                         OBJs as uid
                       } @filter(eq(kind, "device"))
                     }

                     nodes(func: uid(OBJs)) @recurse {
                       children{} 
                       uid
                       name
                       kind
                     }
                   }`

	vars := map[string]string{
		"$account":   request.Account,
		"$namespace": request.Namespace,
	}

	resp, err := txn.QueryWithVars(ctx, q, vars)
	if err != nil {
		return nil, err
	}

	var res struct {
		Nodes []Device `json:"nodes"`
	}

	err = json.Unmarshal(resp.Json, &res)
	if err != nil {
		return nil, err
	}

	var devices []*registrypb.Device
	for _, device := range res.Nodes {
		devices = append(devices, toProto(&device))
	}

	return &registrypb.ListResponse{
		Devices: devices,
	}, nil
}

func (s *Server) ListForAccount(ctx context.Context, request *registrypb.ListDevicesRequest) (response *registrypb.ListResponse, err error) {
	txn := s.dgo.NewReadOnlyTxn()

	const q = `query list($account: string, $namespace: string){
                     var(func: uid($account)) {
                       access.to.namespace @filter(eq(name,$namespace)){
                         owns {
                           OBJs as uid
                         } @filter(not(has(~children)) AND eq(kind, "device"))
                       }
                     }

                     nodes(func: uid(OBJs)) @recurse {
                       children{} 
                       uid
                       name
                       kind
                     }
                   }`

	vars := map[string]string{
		"$account":   request.Account,
		"$namespace": request.Namespace,
	}

	resp, err := txn.QueryWithVars(ctx, q, vars)
	if err != nil {
		return nil, err
	}

	var res struct {
		Nodes []Device `json:"nodes"`
	}

	err = json.Unmarshal(resp.Json, &res)
	if err != nil {
		return nil, err
	}

	var devices []*registrypb.Device
	for _, device := range res.Nodes {
		devices = append(devices, toProto(&device))
	}

	return &registrypb.ListResponse{
		Devices: devices,
	}, nil
}

// TODO make registrypb.Device have multiple certs, ... we also ignore valid_to for now
func toProto(device *Device) *registrypb.Device {
	res := &registrypb.Device{
		Id:      device.UID,
		Name:    device.Name,
		Enabled: &wrappers.BoolValue{Value: device.Enabled},
		Tags:    device.Tags,
		// TODO cert etc
	}

	if len(device.Certificates) > 0 {
		res.Certificate = &registrypb.Certificate{
			PemData:   device.Certificates[0].PemData,
			Algorithm: device.Certificates[0].Algorithm,
		}
	}
	return res
}

func (s *Server) Delete(ctx context.Context, request *registrypb.DeleteRequest) (response *registrypb.DeleteResponse, err error) {
	// // TODO Delete from nodeserver
	// var device Device
	// if err := s.db.First(&device, "name = ?", request.Id).Error; err != nil {
	// 	return nil, err
	// }

	// if err := s.db.Delete(device).Error; err != nil {
	// 	return nil, err
	// }
	return &registrypb.DeleteResponse{}, nil
}
