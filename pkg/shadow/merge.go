package shadow

import (
	"encoding/json"

	"github.com/birdayz/conjungo"
	"github.com/evanphx/json-patch"
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
	var fullJson interface{}
	if full != "" {
		err := json.Unmarshal([]byte(full), &fullJson)
		if err != nil {
			return "", err
		}
	}

	var deltaJson interface{}
	err = json.Unmarshal([]byte(delta), &deltaJson)
	if err != nil {
		return "", err
	}

	err = conjungo.Merge(&fullJson, deltaJson, conjungoOpts)
	if err != nil {
		return "", err
	}

	result, err := json.Marshal(fullJson)
	if err != nil {
		return "", err
	}

	return string(result), nil
}
