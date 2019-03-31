package shadow

import (
	"encoding/json"
	"errors"

	"github.com/birdayz/conjungo"
	jsonpatch "github.com/evanphx/json-patch"
	"github.com/imdario/mergo"
)

func init() {
	conjungoOpts = conjungo.NewOptions()
	conjungoOpts.Overwrite = true
	conjungoOpts.OverwriteDifferentTypes = true

}

var conjungoOpts *conjungo.Options

func calculateDelta(old, new string) string {
	patch, err := jsonpatch.CreateMergePatch([]byte(old), []byte(new))
	if err != nil {
		return new
	}
	return string(patch)
}

func applyDelta(full, delta string) (merged string, err error) {
	if !json.Valid([]byte(delta)) {
		return "", errors.New("delta state is invalid JSON")
	}

	var fullJSON map[string]interface{}
	if full != "" {
		err := json.Unmarshal([]byte(full), &fullJSON)
		if err != nil {
			// full must be a primitive, so we just replace it
			return delta, nil
		}
	}

	var deltaJSON map[string]interface{}
	err = json.Unmarshal([]byte(delta), &deltaJSON)
	if err != nil {
		// delta must be a primitive, so we just replace full with the new primitive value
		return delta, nil
	}

	err = mergo.MergeWithOverwrite(&fullJSON, &deltaJSON)
	if err != nil {
		return "", err
	}

	result, err := json.Marshal(fullJSON)
	if err != nil {
		return "", err
	}

	return string(result), nil
}
