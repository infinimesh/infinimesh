package graph

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/infinimesh/proto/node"
)

func (ctrl *DevicesController) Move(ctx context.Context, msg *connect.Request[node.MoveRequest]) (*connect.Response[node.EmptyMessage], error) {
	req := msg.Msg
	obj := NewBlankDeviceDocument(req.GetUuid())

	return connect.NewResponse(&node.EmptyMessage{}), ctrl.ica_repo.Move(ctx, ctrl, obj, ctrl.ns2dev, req.GetNamespace())
}

func (ctrl *AccountsController) Move(ctx context.Context, _req *connect.Request[node.MoveRequest]) (*connect.Response[node.EmptyMessage], error) {
	req := _req.Msg
	obj := NewBlankAccountDocument(req.GetUuid())

	return connect.NewResponse(&node.EmptyMessage{}), ctrl.ica_repo.Move(ctx, ctrl, obj, ctrl.ns2acc, req.GetNamespace())
}

func StatusFromString(code connect.Code, format string, args ...any) *connect.Error {
	return connect.NewError(code, fmt.Errorf(format, args...))
}
