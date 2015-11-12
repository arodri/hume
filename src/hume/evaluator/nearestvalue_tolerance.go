package evaluator

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"hume/lib/numeric"
	"strconv"
	"sort"
)

type NearestValueTolerance struct {
	BaseEvaluator
	Tolerance
	Value string `json:"value"`
}

func (nv NearestValueTolerance) Evaluate(data map[string]float64, total int) Evaluation {
	var err error
	var msg string
	e:= Evaluation{nv.GetDescription(), "", false}

	testFloat, err := strconv.ParseFloat(nv.Value, 64)
	if err != nil {
		msg = fmt.Sprintf("Error converting value to float: %s", err)
		logrus.Error(msg)
		e.Msg = msg
		return e
	}

	fm, _ := numeric.ND_Mapper(data)
	keys := fm.FloatSlice
	m := fm.Float2String
	sort.Sort(sort.Float64Slice(keys))

	testValue := float64(0)
	for _, k := range keys {
		if testFloat >= k {
			testValue = data[m[k]]
		} else {
			break
		}
	}
	testValue = testValue * 100
	prefix := fmt.Sprintf("Nearest value to %s:%0.2f%%", nv.Value, testValue)
	e = nv.IsOkay(testValue, prefix)
	e.Description = nv.GetDescription()
	return e
}





