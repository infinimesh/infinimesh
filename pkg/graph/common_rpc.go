package graph

import (
	"context"
	"fmt"

	"github.com/arangodb/go-driver"
	"github.com/bufbuild/connect-go"
	inf "github.com/infinimesh/infinimesh/pkg/shared"
	"github.com/infinimesh/proto/node"
	"github.com/infinimesh/proto/node/access"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func _Move(ctx context.Context, c InfinimeshController, obj InfinimeshGraphNode, edge driver.Collection, ns string) error {
	log := c._log().Named("Move")
	log.Debug("Move request received", zap.Any("object", obj), zap.String("namespace", ns))

	requestor := NewBlankAccountDocument(
		ctx.Value(inf.InfinimeshAccountCtxKey).(string),
	)
	log.Debug("Requestor", zap.String("id", requestor.Key))

	err := AccessLevelAndGet(ctx, log, c._DB(), requestor, obj)
	if err != nil {
		return status.Error(codes.NotFound, "Object not found or not enough Access Rights")
	}
	role := obj.GetAccess().Role
	if role != access.Role_OWNER && obj.GetAccess().Level != access.Level_ROOT {
		return status.Error(codes.PermissionDenied, "Must be Owner or Root to perform Move")
	}
	if obj.GetAccess().Namespace == nil {
		return status.Error(codes.Internal, "Object is not under any Namespace, contact support")
	}

	old_namespace := NewBlankNamespaceDocument(*obj.GetAccess().Namespace)

	namespace := NewBlankNamespaceDocument(ns)
	err = AccessLevelAndGet(ctx, log, c._DB(), requestor, namespace)
	if err != nil {
		return status.Error(codes.NotFound, "Namespace not found or not enough Access Rights")
	}
	if namespace.GetAccess().Role != access.Role_OWNER && namespace.GetAccess().Level != access.Level_ROOT {
		return status.Error(codes.PermissionDenied, "Must be Owner or Root to perform Move")
	}

	err = Link(ctx, log, edge, old_namespace, obj, access.Level_NONE, access.Role_UNSET)
	if err != nil {
		log.Warn("Error unlinking Object from Namespace",
			zap.String("object", obj.ID().String()),
			zap.String("namespace", old_namespace.Key),
			zap.Error(err))
		return status.Error(codes.Internal, "Couldn't unlink the object")
	}

	err = Link(ctx, log, edge, namespace, obj, access.Level_ADMIN, role)
	if err != nil {
		log.Warn("Error linking Object to Namespace",
			zap.String("object", obj.ID().String()),
			zap.String("namespace", namespace.Key),
			zap.Error(err))
		return status.Error(codes.Internal, "Couldn't link the object, contact support")
	}

	return nil
}

func (ctrl *DevicesController) Move(ctx context.Context, msg *connect.Request[node.MoveRequest]) (*connect.Response[node.EmptyMessage], error) {
	req := msg.Msg
	obj := NewBlankDeviceDocument(req.GetUuid())

	return connect.NewResponse(&node.EmptyMessage{}), _Move(ctx, ctrl, obj, ctrl.ns2dev, req.GetNamespace())
}

func (ctrl *AccountsController) Move(ctx context.Context, _req *connect.Request[node.MoveRequest]) (*connect.Response[node.EmptyMessage], error) {
	req := _req.Msg
	obj := NewBlankAccountDocument(req.GetUuid())

	return connect.NewResponse(&node.EmptyMessage{}), _Move(ctx, ctrl, obj, ctrl.ns2acc, req.GetNamespace())
}

func StatusFromString(code connect.Code, format string, args ...any) *connect.Error {
	return connect.NewError(code, fmt.Errorf(format, args...))
}
