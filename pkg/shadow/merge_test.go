package shadow

import (
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMergeEmptyString(t *testing.T) {
	old := ""
	merged, err := applyDelta(old, `{"abc" : 13}`)
	require.NoError(t, err)
	require.Equal(t, `{"abc":13}`, merged)
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
	assert.NoError(t, err)

	r := make(map[string]interface{})
	err = json.Unmarshal([]byte(mutatedState), &r)
	assert.NoError(t, err)
	assert.EqualValues(t, 15, r["temp_celsius"])
	assert.EqualValues(t, "20 kmh", r["speed"])
}
