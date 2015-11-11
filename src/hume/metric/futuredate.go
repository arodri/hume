package metric

import (
	"hume/record"
	"time"
)

type FutureDate struct {
	Counter
	Format string `json:"format"`
	Field  string `json:"field"`
	date   time.Time
}

func (fd *FutureDate) Init() error {
	fd.date = time.Now()
	return fd.Counter.Initialize("true","false")
}

func (fd *FutureDate) Process(rec *record.Record) {
	v, ok := rec.Map[fd.Field]
	l := "false"
	if ok {
		d, err := time.Parse(fd.Format, v)
		if err == nil && d.After(fd.date) {
			l = "true"
		}
	}
	fd.Count(l)
}
