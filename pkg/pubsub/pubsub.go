package pubsub

import (
	_pubsub "github.com/cskr/pubsub"
)

type PubSub interface {
	Sub(topics ...string) chan interface{}
	SubOnce(topics ...string) chan interface{}
	SubOnceEach(topics ...string) chan interface{}
	AddSub(ch chan interface{}, topics ...string)
	AddSubOnceEach(ch chan interface{}, topics ...string)

	Pub(msg interface{}, topics ...string)
	TryPub(msg interface{}, topics ...string)

	Unsub(ch chan interface{}, topics ...string)
	Close(topics ...string)
	Shutdown()
}

func New(capacity int) PubSub {
	return _pubsub.New(capacity)
}
