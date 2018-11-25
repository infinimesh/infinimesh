package registry

import (
	"errors"
	"fmt"
	"log"

	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"

	"encoding/base64"

	"github.com/google/uuid"
	"github.com/infinimesh/infinimesh/pkg/registry/registrypb"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // nolint: golint
	context "golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	if request.Certificate == nil {
		return nil, status.Error(codes.FailedPrecondition, "No certificate provided")
	}
	st, err := base64.StdEncoding.DecodeString(request.Certificate.PemData)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, "PEM data is not valid base64")
	}
	fp, err := s.getFingerprint(st, request.Certificate.Algorithm)
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
		Namespace:              request.Namespace,
		Name:                   request.Name,
		Enabled:                request.Enabled,
		Certificate:            string(st),
		CertificateType:        request.Certificate.Algorithm,
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
		Name:      device.Name,
		Namespace: device.Namespace,
	}, nil
}
