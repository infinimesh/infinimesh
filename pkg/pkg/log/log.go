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

package log

import (
	"io/ioutil"
	"strings"

	"go.uber.org/zap"
)

// NewProdOrDev returns a Prod logger if the application is running in Docker,
// otherwise it returns a Dev logger.
func NewProdOrDev() (log *zap.Logger, err error) {
	if runningInDocker() {
		return zap.NewProduction()
	} else {
		return zap.NewDevelopment()
	}
}

// TODO Will this work with CRI-O, rkt and other container runtimes beside docker?
func runningInDocker() bool {
	f, err := ioutil.ReadFile("/proc/1/cgroup")
	if err != nil {
		return false
	}

	if strings.Contains(string(f), "docker") {
		return true
	}
	return false
}
