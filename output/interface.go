package output

import (
	"FuckEventvwr/config"
	"FuckEventvwr/velocidex/evtx"
)

type output interface {
	WriteRecord(record *evtx.EventRecord) error
	WriteError(err string) error
	Close() error
}

var Output output

func InitOutput() {
	if config.Cfg.EsURL != "" {
		Output = NewElasticsearch()
	} else {
		Output = NewExcel()
	}
}
