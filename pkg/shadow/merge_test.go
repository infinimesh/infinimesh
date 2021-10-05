//--------------------------------------------------------------------------
// Copyright 2018 infinimesh
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
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/require"
)

func TestMergeEmptyString(t *testing.T) {
	old := ""
	merged, err := applyDelta(old, `{"abc" : 13}`)
	require.NoError(t, err)
	require.JSONEq(t, `{"abc":13}`, merged)
}

func TestMerge(t *testing.T) {
	state := `
{
"temp_celsius" : 20,
"speed" : "20 kmh"
}
`

	delta := `
{
"temp_celsius" : 15
}
`

	mutatedState, err := applyDelta(state, delta)
	require.NoError(t, err)

	r := make(map[string]interface{})
	err = json.Unmarshal([]byte(mutatedState), &r)
	require.NoError(t, err)
	require.EqualValues(t, 15, r["temp_celsius"])
	require.EqualValues(t, "20 kmh", r["speed"])
}

func TestMergeNested(t *testing.T) {
	state := `{"a":false}`

	delta := `{"a":{"very_much":true,"bla":13}}`
	mutatedState, err := applyDelta(state, delta)
	require.NoError(t, err)

	require.NoError(t, err)
	require.EqualValues(t, `{"a":{"bla":13,"very_much":true}}`, mutatedState)

}

func TestMergePrimitive(t *testing.T) {
	old := `{"a" : "test"}`
	newDelta := `true`

	merged, err := applyDelta(old, newDelta)

	require.NoError(t, err)
	require.EqualValues(t, "true", merged)
}

func TestMergeBothPrimitive(t *testing.T) {
	old := `true`
	newDelta := `false`

	merged, err := applyDelta(old, newDelta)

	require.NoError(t, err)
	require.JSONEq(t, "false", merged)
}

func TestMergeOldPrimitive(t *testing.T) {
	old := `true`
	newDelta := `{"a" : "b"}`

	merged, err := applyDelta(old, newDelta)

	require.NoError(t, err)
	require.JSONEq(t, newDelta, merged)
}

func TestCalculateDelta(t *testing.T) {
	old := `{"a":{"very_much":true}}`
	new := `{"a":{"very_much":true,"bla":13}}`

	patch := calculateDelta(old, new)

	expected := `{"a":{"bla":13}}`
	require.EqualValues(t, expected, patch)
}

func TestCalculateDeltaArray(t *testing.T) {
	full := `{"a":["abc","def"]}`
	new := `{"a":["fitze","fatze"]}`

	merged, err := applyDelta(full, new)
	require.NoError(t, err)

	expected := `{"a":["fitze","fatze"]}`
	require.EqualValues(t, expected, merged)
}

// Map string -> primitive is replaced by string -> object
func TestNestedDifferentType(t *testing.T) {
	old := `{"2_202":21.31999969482422,"2_203":823.4000244140625}`
	new := `{"2_202":{"value":21.639999389648438,"name":"Measurand Room Temperature"},"2_203":{"value":823.4000244140625,"name":"Measurand Volatile Organic Components for Air Quality"}}`

	merged, err := applyDelta(old, new)

	require.NoError(t, err)
	require.JSONEq(t, new, merged)
}
