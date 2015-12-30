package metric

import (
	"hume/record"
)

type NumRecords struct {
	BaseMetric
	Counter
}

func (nr *NumRecords) Init() error {
	return nr.Counter.Initialize()
}

func (nr *NumRecords) Process(rec *record.Record) {
	nr.Count("true")
}
