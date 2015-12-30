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

func (nv *NearestValueTolerance) compute(data map[string]float64, total int) (float64, error, string) {
	var err error
	var msg string
	testValue := float64(0)

	testFloat, err := strconv.ParseFloat(nv.Value, 64)
	if err != nil {
		msg = fmt.Sprintf("Error converting value to float: %s", err)
		logrus.Error(msg)
		return testValue, err, msg
	}

	fm, _ := numeric.ND_Mapper(data)
	keys := fm.FloatSlice
	m := fm.Float2String
	sort.Sort(sort.Float64Slice(keys))

	for _, k := range keys {
		if testFloat >= k {
			testValue = data[m[k]]
		} else {
			break
		}
	}

	testValue = testValue * 100
	return testValue, err, msg

}

func (nv *NearestValueTolerance) Evaluate(data map[string]float64, total int) Evaluation {
	e:= Evaluation{nv.GetDescription(), "", false}
	testValue, err, msg := nv.compute(data, total)
	if err != nil {
		e.Msg = msg
		return e
	}

	prefix := fmt.Sprintf("Nearest value to %s:%0.2f%%", nv.Value, testValue)
	e = nv.IsOkay(testValue, prefix)
	e.Description = nv.GetDescription()
	return e
}

func (nv *NearestValueTolerance) Train(data map[string]float64, total int) error {
	testValue, err, _ := nv.compute(data, total)
	if err != nil {
		return err
	}
	nv.SetValue(testValue)
	return nil
}
