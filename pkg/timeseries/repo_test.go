package timeseries

import (
	"context"
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
}

func TestSave(t *testing.T) {
	err := repo.CreateDataPoint(context.TODO(), &DataPoint{
		DeviceID:   "test-device-1",
		DeviceName: "test-device-1",
		Property:   "voltage",
		Timestamp:  time.Now(),
		Value:      50.0,
	})
	require.NoError(t, err)
}
