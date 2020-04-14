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
