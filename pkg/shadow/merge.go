package shadow

import (
	"encoding/json"
	"errors"

	"github.com/birdayz/conjungo"
	jsonpatch "github.com/evanphx/json-patch"
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

	res, err := jsonpatch.MergePatch([]byte(full), []byte(delta))
	if err != nil {
		return delta, nil
	}

	return string(res), nil
}
