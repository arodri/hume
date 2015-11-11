package metric

import (
	"hume/record"
)

type ValidValues struct {
	Counter
	Field    string   `json:"field"`
	Values   []string `json:"values"`
	valueMap map[string]string
}

func (vv *ValidValues) Init() error {
	vv.valueMap = make(map[string]string)
	for _, val := range vv.Values {
		vv.valueMap[val] = val
	}
	return vv.Counter.Initialize("true","false")
}

func (vv *ValidValues) Process(rec *record.Record) {
	v, _ := rec.Map[vv.Field]
	_, ok := vv.valueMap[v]
	l := "false"
	if ok {
		l = "true"
	}
	vv.Count(l)
}
