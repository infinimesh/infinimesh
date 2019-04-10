package timeseries

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type TimeseriesRepo interface {
	CreateDataPoint(ctx context.Context, datapoint *DataPoint) error
}

type DataPoint struct {
	DeviceID   string
	DeviceName string
	Property   string
	Timestamp  time.Time
	Value      float64
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

	_, err = tx.Exec("INSERT INTO DATA_POINTS (device_id, device_name, property, timestamp, value) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING",
		datapoint.DeviceID, datapoint.DeviceName, datapoint.Property, datapoint.Timestamp, datapoint.Value,
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
