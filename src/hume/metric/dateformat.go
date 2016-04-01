package metric

import (
	"hume/record"
	"time"
)

type DateFormat struct {
	BaseMetric
	Counter
	Format string `json:"format"`
	Field  string `json:"field"`
}

func (df *DateFormat) Init() error {
	return df.Counter.Initialize("true", "false")
}

func (df *DateFormat) Process(rec *record.Record) {
	v, ok := rec.Map[df.Field]
	l := "false"
	if ok {
		_, err := time.Parse(df.Format, v)
		if err == nil || v == "" {
			l = "true"
		}
	}
	df.Count(l)
}
