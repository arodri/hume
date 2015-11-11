package metric

import (
	"hume/record"
)

type Populated struct {
	Counter
	Field string `json:"field"`
}

func (p *Populated) Init() error {
	return p.Counter.Initialize("true","false")
}

func (p *Populated) Process(rec *record.Record) {
	v, ok := rec.Map[p.Field]
	l := "true"
	if !ok || v == "" {
		l = "false"
	}
	p.Count(l)
}
