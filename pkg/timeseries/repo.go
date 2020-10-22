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
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

const (
	sizeKB = 1 << (10 * 1)
)

type TimeseriesRepo interface {
	CreateDataPoint(ctx context.Context, datapoint *DataPoint) error
	ReadExistingDatapoint(ctx context.Context, deviceID string) (float32, error)
}

type DataPoint struct {
	DeviceID  string
	MessageID uint64
	Property  string
	Timestamp time.Time
	Value     float64
	Length    float32
}

type timescaleRepo struct {
	log *zap.Logger
	db  *sql.DB
}

func NewTimescaleRepo(log *zap.Logger, connection string) (result TimeseriesRepo, err error) {
	conn, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, err
	}

	// set connection limit -> https://godoc.org/database/sql#DB.SetMaxOpenConns
	conn.SetMaxIdleConns(0)
	conn.SetMaxOpenConns(90)

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	return &timescaleRepo{
		log: log,
		db:  conn,
	}, nil
}

func (t *timescaleRepo) CreateDataPoint(ctx context.Context, datapoint *DataPoint) error {
	tx, err := t.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO DATA_POINTS (device_id, message_id, property, timestamp, value, message_length) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING",
		datapoint.DeviceID, datapoint.MessageID, datapoint.Property, datapoint.Timestamp, datapoint.Value, datapoint.Length,
	)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (t *timescaleRepo) ReadExistingDatapoint(ctx context.Context, deviceID string) (float32, error) {
	tx, err := t.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, err
	}
	row := tx.QueryRow("SELECT message_length FROM DATA_POINTS where device_id= $1 ORDER BY timestamp DESC LIMIT 1", deviceID)
	var messageLength float32
	err = row.Scan(&messageLength)
	if err != nil {
		fmt.Printf("no existing rows %v", err)
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return messageLength, nil
}
