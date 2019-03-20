package shadow

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSetGet(t *testing.T) {
	key := uuid.New().String()
	repo, err := NewRedisRepo(":6379")
	require.NoError(t, err)

	input := DeviceState{
		ID:      key,
		Version: 1,
		State:   json.RawMessage([]byte("50")),
	}

	err = repo.SetDesired(input)
	require.NoError(t, err)

	ds, err := repo.GetDesired(key)

	require.NoError(t, err)
	require.EqualValues(t, input, ds)
}

func TestSetGetDesiredAndReported(t *testing.T) {
	key := uuid.New().String()
	repo, err := NewRedisRepo(":6379")
	require.NoError(t, err)

	input := DeviceState{
		ID:      key,
		Version: 1,
		State:   json.RawMessage([]byte("50")),
	}

	err = repo.SetDesired(input)
	require.NoError(t, err)

	inputReported := DeviceState{
		ID:      key,
		Version: 1,
		State:   json.RawMessage([]byte("60")),
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
