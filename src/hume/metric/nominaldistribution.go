package metric

import(
	"hume/record"
)

type NominalDistribution struct {
	Counter
	Field string `json:"field"`

}

func (nd *NominalDistribution) Process(rec *record.Record) {
	v, _ := rec.Map[nd.Field]
	nd.Count(v)
}
