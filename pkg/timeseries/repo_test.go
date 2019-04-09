package timeseries

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	log *zap.Logger
)

func init() {
	log, _ = zap.NewDevelopment()
}

func TestSave(t *testing.T) {
	repo, err := NewTimescaleRepo(log, "postgres://postgres:postgres@localhost/postgres?sslmode=disable")
	require.NoError(t, err)
	err = repo.CreateDataPoint(context.TODO(), &DataPoint{
		DeviceID:   "test-device-1",
		DeviceName: "test-device-1",
		Property:   "voltage",
		Timestamp:  time.Now(),
		Value:      50.0,
	})
	require.NoError(t, err)
}
