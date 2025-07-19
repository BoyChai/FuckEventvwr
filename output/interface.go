package output

import (
	"FuckEventvwr/velocidex/evtx"
)

type output interface {
	WriteRecord(record *evtx.EventRecord) error
	WriteError(err string) error
	Count() (int, error)
	Close() (int, error)
}

var Module output
