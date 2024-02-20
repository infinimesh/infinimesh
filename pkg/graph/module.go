package graph

import (
	"context"

	"github.com/arangodb/go-driver"
	"github.com/go-redis/redis/v8"
	"github.com/infinimesh/infinimesh/pkg/credentials"
	"github.com/infinimesh/infinimesh/pkg/sessions"
	"github.com/infinimesh/proto/handsfree"
	"github.com/infinimesh/proto/node/accounts"
	"github.com/infinimesh/proto/node/devices"
	"github.com/infinimesh/proto/node/nodeconnect"
	"go.uber.org/zap"
)

type DevicesControllerModule interface {
	Handler() nodeconnect.DevicesServiceHandler
	SetSigningKey([]byte)
}

type devicesControllerModule struct {
	handler *DevicesController
}

func (m *devicesControllerModule) Handler() nodeconnect.DevicesServiceHandler {
	return m.handler
}

func (m *devicesControllerModule) SetSigningKey(key []byte) {
	m.handler.SIGNING_KEY = key
}

func NewDevicesControllerModule(log *zap.Logger, db driver.Database,
	hfc handsfree.HandsfreeServiceClient, bus *EventBus) DevicesControllerModule {
	return &devicesControllerModule{
		handler: NewDevicesController(
			log, db, hfc,
			NewInfinimeshCommonActionsRepo(log.Named("DevicesController"), db),
			NewGenericRepo[*devices.Device](db),
			NewGenericRepo[*accounts.Account](db),
			bus,
		),
	}
}

type AccountsControllerModule interface {
	Handler() nodeconnect.AccountsServiceHandler
	SetSigningKey([]byte)
}

type accountsControllerModule struct {
	handler *AccountsController
}

func (m *accountsControllerModule) Handler() nodeconnect.AccountsServiceHandler {
	return m.handler
}

func (m *accountsControllerModule) SetSigningKey(key []byte) {
	m.handler.SIGNING_KEY = key
}

func NewAccountsControllerModule(log *zap.Logger, db driver.Database, rdb redis.Cmdable, bus *EventBus) AccountsControllerModule {
	return &accountsControllerModule{
		handler: NewAccountsController(
			log, db, rdb,
			sessions.NewSessionsHandlerModule(rdb).Handler(),
			NewInfinimeshCommonActionsRepo(log.Named("AccountsController"), db),
			NewGenericRepo[*accounts.Account](db),
			credentials.NewCredentialsController(context.Background(), log, db),
			bus,
		),
	}
}
