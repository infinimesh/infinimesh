/*
Copyright © 2021-2022 Infinite Devices GmbH, Nikita Ivanovski info@slnt-opp.xyz

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package graph

import (
	"context"

	"github.com/arangodb/go-driver"
	"github.com/infinimesh/infinimesh/pkg/graph/schema"
	pb "github.com/infinimesh/infinimesh/pkg/node/proto"
	devpb "github.com/infinimesh/infinimesh/pkg/node/proto/devices"
	"go.uber.org/zap"
)

type Device struct {
	*devpb.Device
	driver.DocumentMeta	
}

func (o *Device) ID() (driver.DocumentID) {
	return o.DocumentMeta.ID
}

func NewBlankDeviceDocument(key string) (*Device) {
	return &Device{
		Device: &devpb.Device{
			Uuid: key,
		},
		DocumentMeta: NewBlankDocument(schema.DEVICES_COL, key),
	}
}


type DevicesController struct {
	pb.UnimplementedDevicesServiceServer
	log *zap.Logger

	col driver.Collection // Devices Collection
	db driver.Database

	ns2dev driver.Collection // Namespaces to Devices permissions edge collection

	SIGNING_KEY []byte
}

func NewDevicesController(log *zap.Logger, db driver.Database) DevicesController {
	ctx := context.TODO()
	col, _ := db.Collection(ctx, schema.DEVICES_COL)

	return DevicesController{
		log: log.Named("DevicesController"), col: col, db: db,
		ns2dev: GetEdgeCol(ctx, db, schema.NS2DEV),
		SIGNING_KEY: []byte("just-an-init-thing-replace-me"),
	}
}