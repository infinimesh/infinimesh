package shadow

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var repo Repo

func init() {
	dbURL := os.Getenv("DB_ADDR")
	if dbURL == "" {
		dbURL = ":6379"
	}
	r, err := NewRedisRepo(dbURL)
	if err != nil {
		panic(err)
	}
	repo = r
}

func TestSetGet(t *testing.T) {
	key := uuid.New().String()

	input := DeviceState{
		ID: key,
		State: DeviceStateMessage{
			Version: 1,
			State:   json.RawMessage([]byte("50")),
		},
	}

	err := repo.SetDesired(input)
	require.NoError(t, err)

	ds, err := repo.GetDesired(key)

	require.NoError(t, err)
	require.EqualValues(t, input, ds)
}

func TestSetGetDesiredAndReported(t *testing.T) {
	key := uuid.New().String()

	input := DeviceState{
		ID: key,
		State: DeviceStateMessage{
			Version: 1,
			State:   json.RawMessage([]byte("50")),
		},
	}

	err := repo.SetDesired(input)
	require.NoError(t, err)

	inputReported := DeviceState{
		ID: key,
		State: DeviceStateMessage{
			Version: 1,
			State:   json.RawMessage([]byte("60")),
		},
	}

	err = repo.SetReported(inputReported)
	require.NoError(t, err)

	ds, err := repo.GetDesired(key)
	require.NoError(t, err)
	require.EqualValues(t, input, ds)

	rs, err := repo.GetReported(key)
	require.NoError(t, err)
	require.EqualValues(t, inputReported, rs)

}
