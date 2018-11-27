package shadow

import (
	"encoding/json"

	"github.com/birdayz/conjungo"
)

func init() {
	conjungoOpts = conjungo.NewOptions()
	conjungoOpts.Overwrite = true
	conjungoOpts.OverwriteDifferentTypes = true

}

var conjungoOpts *conjungo.Options

func applyDelta(full, delta string) (merged string, err error) {
	fullMap := make(map[string]interface{})
	if full != "" {
		err := json.Unmarshal([]byte(full), &fullMap)
		if err != nil {
			return "", err
		}
	}

	deltaMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(delta), &deltaMap)
	if err != nil {
		return "", err
	}

	err = conjungo.Merge(&fullMap, deltaMap, conjungoOpts)
	if err != nil {
		return "", err
	}

	result, err := json.Marshal(fullMap)
	if err != nil {
		return "", err
	}

	return string(result), nil
}
