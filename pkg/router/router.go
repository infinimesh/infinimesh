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

package router

import (
	"fmt"
	"strings"
)

type Router struct {
	fallback string
}

func New(fallback string) *Router {
	return &Router{
		fallback: fallback,
	}
}

func (r *Router) Route(inputTopic, inputDevice string) (outputTopic string) {
	splt := strings.Split(inputTopic, "/")
	switch splt[0] {
	case "devices":
		// Check if at least the segment for the deviceID is given plus at least one subtopic segment
		if len(splt) >= 3 {
			deviceInTopic := splt[1]
			if inputDevice != deviceInTopic && deviceInTopic != "+" {
				fmt.Println("Input topic does not match device.", deviceInTopic, inputDevice)
			} else {
				subtopic := strings.Join(splt[2:], "/")
				fmt.Println("Subtopic:", subtopic)
				switch subtopic {
				case "state/reported/delta":
					return "shadow.reported-state.delta"
				}
			}
		}
	}

	return r.fallback
}
