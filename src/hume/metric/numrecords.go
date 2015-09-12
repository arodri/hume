package metric

import (
	"hume/record"
)

type NumRecords struct {
	Counter
}

func (nr *NumRecords) Process(rec *record.Record) {
	nr.Count("true")
}
