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
