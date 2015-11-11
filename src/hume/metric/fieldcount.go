package metric

import (
	//	"github.com/Sirupsen/logrus"
	// "sync"
	"hume/record"
	"strconv"
)

type FieldCount struct {
	Counter
}

func (ft *FieldCount) Init() error {
	return ft.Counter.Initialize()
}

func (ft *FieldCount) Process(rec *record.Record) {
	l := strconv.Itoa(rec.ValueCount())
	ft.Count(l)
}
