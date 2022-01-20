package repo

import (
	"context"

	"github.com/slntopp/infinimesh/pkg/repo/repopb"
	"go.uber.org/zap"
)

type Server struct {
	Repo Repo
	Log  *zap.Logger
}

func (s *Server) Get(context context.Context, request *repopb.GetRequest) (response *repopb.GetResponse, err error) {
	log := s.Log.Named("Get Device Status Controller")
	log.Info("Function Invoked", zap.String("Device", request.Id))
	response = &repopb.GetResponse{
		Repo: &repopb.Repo{},
	}
	deviceStatus, err := s.Repo.GetDeviceStatus(request.Id)
	if err != nil {
		deviceStatus.ID = request.Id
		deviceStatus.Status = DeviceStatus{
			Fingerprint: []byte("{}"),
			Enabled:     false,
			NamespaceID: "",
		}
	}
	response.Repo = &repopb.Repo{
		Enabled:     	deviceStatus.Status.Enabled,
		BasicEnabled: deviceStatus.Status.BasicEnabled,
		NamespaceID: 	deviceStatus.Status.NamespaceID,
		FingerPrint: 	deviceStatus.Status.Fingerprint,
	}
	return response, nil
}

func (s *Server) SetDeviceState(context context.Context, request *repopb.SetDeviceStateRequest) (response *repopb.SetDeviceStateResponse, err error) {
	log := s.Log.Named("Set Device Status Controller")
	log.Info("Function Invoked", zap.String("Device", request.Id))
	deviceState := DeviceState{
		ID: request.Id,
		Status: DeviceStatus{
			Fingerprint: request.Repo.FingerPrint,
			Enabled:     request.Repo.Enabled,
			NamespaceID: request.Repo.NamespaceID,
		},
	}
	err = s.Repo.SetDeviceStatus(deviceState)
	if err != nil {
		response := &repopb.SetDeviceStateResponse{
			Status: false,
		}
		return response, err
	}
	response = &repopb.SetDeviceStateResponse{
		Status: true,
	}
	return response, nil
}

func (s *Server) DeleteDeviceState(context context.Context, request *repopb.DeleteDeviceStateRequest) (response *repopb.DeleteDeviceStateResponse, err error) {
	log := s.Log.Named("Delete Device Status Controller")
	log.Info("Function Invoked", zap.String("Device", request.Id))
	err = s.Repo.DeleteDeviceStatus(request.Id)
	if err != nil {
		response := &repopb.DeleteDeviceStateResponse{
			Status: false,
		}
		return response, err
	}
	response = &repopb.DeleteDeviceStateResponse{
		Status: true,
	}
	return response, nil
}
