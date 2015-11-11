package metric

import (
	"hume/record"
	"time"
)

type DateFormat struct {
	Counter
	Format string `json:"format"`
	Field  string `json:"field"`
}

func (df *DateFormat) Init() error {
	return df.Counter.Initialize("true","false")
}

func (df *DateFormat) Process(rec *record.Record) {
	v, ok := rec.Map[df.Field]
	l := "true"
	if ok {
		_, err := time.Parse(df.Format, v)
		if err != nil {
			l = "false"
		}
	}
	df.Count(l)
}
