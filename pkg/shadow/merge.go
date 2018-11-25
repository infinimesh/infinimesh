package shadow

import (
	"encoding/json"

	"github.com/imdario/mergo"
)

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

	err = mergo.Merge(&fullMap, deltaMap, mergo.WithOverride)
	if err != nil {
		return "", err
	}

	result, err := json.Marshal(fullMap)
	if err != nil {
		return "", err
	}

	return string(result), nil
}
