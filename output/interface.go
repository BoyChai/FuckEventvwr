package output

import "FuckEventvwr/velocidex/evtx"

type output interface {
	Write(record *evtx.EventRecord) error
	Close() error
}

var Output output
