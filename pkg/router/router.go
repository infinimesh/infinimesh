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
	fmt.Println("r")
	splt := strings.Split(inputTopic, "/")
	switch splt[0] {
	case "devices":
		// Check if at least the segment for the deviceID is given plus at least one subtopic segment
		if len(splt) >= 3 {
			deviceInTopic := splt[1]
			if inputDevice != deviceInTopic {
				// TODO Currently, reject these. However, this will
				// not be an error anymore once devices can be
				// authorized to send on behalf of other devices
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
