//--------------------------------------------------------------------------
// Copyright 2018 Infinite Devices GmbH
// www.infinimesh.io
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//--------------------------------------------------------------------------

package timeseries

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	log  *zap.Logger
	repo TimeseriesRepo
)

func init() {
	log, _ = zap.NewDevelopment()

	dbAddr := os.Getenv("TIMESERIES_URL")
	if dbAddr == "" {
		dbAddr = "postgres://postgres:postgres@localhost/postgres?sslmode=disable"
	}

	r, err := NewTimescaleRepo(log, dbAddr)
	if err != nil {
		panic(err)
	}
	repo = r

	conn, err := sql.Open("postgres", dbAddr)
	if err != nil {
		panic(err)
	}
	ddl, err := ioutil.ReadFile("./ddl/ddl.sql")
	if err != nil {
		panic(err)
	}

	txn, err := conn.Begin()
	if err != nil {
		panic(err)
	}
	_, err = txn.Exec(string(ddl))
	if err != nil {
		fmt.Printf("Error during DDL import: %v. Ignoring\n", err)
	}
	txn.Commit()

}

func TestSave(t *testing.T) {
	err := repo.CreateDataPoint(context.TODO(), &DataPoint{
		DeviceID:  "test-device-1",
		MessageID: uint64(2),
		Property:  "voltage",
		Timestamp: time.Now(),
		Value:     50.0,
		Length:    12.0,
	})
	require.NoError(t, err)
}

func TestRead(t *testing.T) {
	messageLength, err := repo.ReadExistingDatapoint(context.TODO(), "test-device-1", 2)
	err = repo.CreateDataPoint(context.TODO(), &DataPoint{
		DeviceID:  "test-device-1",
		MessageID: uint64(2),
		Property:  "voltage",
		Timestamp: time.Now(),
		Value:     50.0,
		Length:    12.0 + messageLength,
	})
	require.NoError(t, err)
}
