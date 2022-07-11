/*
Copyright © 2021-2022 Infinite Devices GmbH Nikita Ivanovski info@slnt-opp.xyz

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package handsfree

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"strconv"
	"time"
)

func GenerateCodeLong(app string) (string, error) {
	s := app + strconv.FormatInt(time.Now().UnixMicro(), 10)
	hasher := md5.New()

	_, err := hasher.Write([]byte(s))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func ShortenToFit[T any](code string, db map[string]T) (string, error) {
	if len(code) != 32 {
		return "", errors.New("code isn't md5 hash(must be 32 char)")
	}

	for i := 0; i < 29; i++ {
		r := code[i:i+3] + code[29-i:32-i]
		if _, ok := db[r]; !ok {
			return r, nil
		}
	}

	return "", errors.New("couldn't find the fitting code")
}
