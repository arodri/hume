package metric

import(
	"hume/lib/numeric"
	"hume/record"
	"strconv"
)

type NumericDistribution struct {
	Counter
	Field string `json:"field"`
}

func (nd *NumericDistribution) Init() error {
	return nd.Counter.Initialize()
}

func (nd *NumericDistribution) Process(rec *record.Record) {
	var v string
	s, _ := rec.Map[nd.Field]
	if s == "" {
		v = numeric.EMPTY_STRING
	} else {
		_, err := strconv.ParseFloat(s, 64)
		if err != nil {
			v = numeric.NOT_A_FLOAT
		} else {
			v = s
		}
	}
	nd.Count(v)
}
