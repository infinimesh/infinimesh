package graph

import (
	"context"
	"fmt"
	proto_eventbus "github.com/infinimesh/proto/eventbus"
	"go.uber.org/zap"

	"connectrpc.com/connect"
	"github.com/infinimesh/proto/node"
)

func (ctrl *DevicesController) Move(ctx context.Context, msg *connect.Request[node.MoveRequest]) (*connect.Response[node.EmptyMessage], error) {
	log := ctrl.Log().Named("MoveDevice")

	req := msg.Msg
	obj := NewBlankDeviceDocument(req.GetUuid())

	err := ctrl.ica_repo.Move(ctx, ctrl, obj, ctrl.ns2dev, req.GetNamespace())
	if err != nil {
		log.Error("Failed to move device", zap.Error(err))
		return nil, err
	}

	err = ctrl.bus.Notify(ctx, &proto_eventbus.Event{
		EventKind: proto_eventbus.EventKind_DEVICE_MOVE,
		Entity:    &proto_eventbus.Event_Device{Device: obj.Device},
	})

	if err != nil {
		log.Error("Failed to notify move", zap.Error(err))
	}

	return connect.NewResponse(&node.EmptyMessage{}), nil
}

func (ctrl *AccountsController) Move(ctx context.Context, _req *connect.Request[node.MoveRequest]) (*connect.Response[node.EmptyMessage], error) {
	log := ctrl.Log().Named("MoveAccount")

	req := _req.Msg
	obj := NewBlankAccountDocument(req.GetUuid())

	err := ctrl.ica_repo.Move(ctx, ctrl, obj, ctrl.ns2acc, req.GetNamespace())
	if err != nil {
		log.Error("Failed to move device", zap.Error(err))
		return nil, err
	}

	err = ctrl.bus.Notify(ctx, &proto_eventbus.Event{
		EventKind: proto_eventbus.EventKind_ACCOUNT_MOVE,
		Entity:    &proto_eventbus.Event_Account{Account: obj.Account},
	})

	if err != nil {
		log.Error("Failed to notify move", zap.Error(err))
	}

	return connect.NewResponse(&node.EmptyMessage{}), ctrl.ica_repo.Move(ctx, ctrl, obj, ctrl.ns2acc, req.GetNamespace())
}

func StatusFromString(code connect.Code, format string, args ...any) *connect.Error {
	return connect.NewError(code, fmt.Errorf(format, args...))
}
