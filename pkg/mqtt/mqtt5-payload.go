package mqtt

import (
	"time"
)

type topic struct {
	name string `json:name`
	data []byte `json:byte[]`
}

type message struct {
	topics []topic `json:"topics"`
}

type payload struct {
	version   string    `json:version`
	timestamp time.Time `json:timestamp`
	message   message   `json:message`
}
