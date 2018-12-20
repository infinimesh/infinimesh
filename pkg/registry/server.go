package registry

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/infinimesh/infinimesh/pkg/registry/registrypb"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Server struct {
	db *gorm.DB
}

func NewServer(addr string) *Server {
	db, err := gorm.Open("postgres", addr)
	if err != nil {
		log.Fatal(err)
	}

	db.LogMode(false)
	db.SingularTable(true)
	db.AutoMigrate(&Device{})

	return &Server{
		db: db,
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
	if request.Device.Certificate == nil {
		return nil, status.Error(codes.FailedPrecondition, "No certificate provided")
	}
	st, err := base64.StdEncoding.DecodeString(request.Device.Certificate.PemData)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, "PEM data is not valid base64")
	}
	fp, err := s.getFingerprint(st, request.Device.Certificate.Algorithm)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, "Invalid Certificate")
	}
	u, err := uuid.NewRandom()
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}
	uuidBytes, err := u.MarshalBinary()
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	if err := s.db.Create(&Device{
		ID:                     uuidBytes,
		Name:                   request.Device.Id,
		Enabled:                request.Device.Enabled,
		Certificate:            string(st),
		CertificateType:        request.Device.Certificate.Algorithm,
		CertificateFingerprint: fp,
	}).Error; err != nil {
		return nil, status.Error(codes.FailedPrecondition, fmt.Sprintf("Failed to create device: %v", err))
	}
	return &registrypb.CreateResponse{
		Fingerprint: fp,
	}, nil
}

func (s *Server) GetByFingerprint(ctx context.Context, request *registrypb.GetByFingerprintRequest) (*registrypb.GetByFingerprintResponse, error) {
	device := &Device{}
	if err := s.db.Take(device, &Device{CertificateFingerprint: request.Fingerprint}).Error; err != nil {
		return &registrypb.GetByFingerprintResponse{}, status.Error(codes.FailedPrecondition, err.Error())
	}
	return &registrypb.GetByFingerprintResponse{
		Id: device.Name,
	}, nil
}

func (s *Server) Get(ctx context.Context, request *registrypb.GetRequest) (response *registrypb.GetResponse, err error) {
	var device Device
	if err := s.db.First(&device, "name = ?", request.Id).Error; err != nil {
		return nil, err
	}
	return &registrypb.GetResponse{
		Id:      device.Name,
		Enabled: false,
		// TODO
	}, nil
}
func (s *Server) List(context.Context, *registrypb.ListDevicesRequest) (*registrypb.ListResponse, error) {
	var devices []*Device
	if err := s.db.Find(&devices).Error; err != nil {
		return nil, err
	}

	var protoDevices []*registrypb.Device
	for _, device := range devices {
		protoDevices = append(protoDevices, toProto(device))
	}

	return &registrypb.ListResponse{
		Devices: protoDevices,
	}, nil
}

func toProto(device *Device) *registrypb.Device {
	return &registrypb.Device{
		Id:      device.Name,
		Enabled: true, // TODO
		Certificate: &registrypb.Certificate{
			PemData:   device.Certificate,
			Algorithm: "",        // TODO
			ValidTo:   uint64(0), // TODO
		},
	}
}

func (s *Server) Delete(ctx context.Context, request *registrypb.DeleteRequest) (response *registrypb.DeleteResponse, err error) {
	if err := s.db.Delete(&Device{Name: request.Id}).Error; err != nil {
		return nil, err
	}
	return &registrypb.DeleteResponse{}, nil
}
