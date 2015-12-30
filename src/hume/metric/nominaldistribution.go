package metric

import (
	"hume/record"
)

type NominalDistribution struct {
	BaseMetric
	Counter
	Field string `json:"field"`
}

func (nd *NominalDistribution) Init() error {
	return nd.Counter.Initialize()
}

func (nd *NominalDistribution) Process(rec *record.Record) {
	v, _ := rec.Map[nd.Field]
	nd.Count(v)
}
