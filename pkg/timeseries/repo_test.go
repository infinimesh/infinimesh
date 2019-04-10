package timeseries

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"io/ioutil"

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
		DeviceID:   "test-device-1",
		DeviceName: "test-device-1",
		Property:   "voltage",
		Timestamp:  time.Now(),
		Value:      50.0,
	})
	require.NoError(t, err)
}
