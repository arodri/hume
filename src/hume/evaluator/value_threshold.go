package evaluator

import (
	"fmt"
	"github.com/Sirupsen/logrus"
)

type ValueThreshold struct {
	BaseEvaluator
	Thresholder
	Value    string `json:"value"`
	UseCount bool   `json:"use_count"`
}

func (v ValueThreshold) Evaluate(data map[string]float64, total int) Evaluation {
	cnt, ok := data[v.Value]
	if !ok {
		logrus.Error(fmt.Sprintf("Value %#v is not valid", v.Value))
	}

	testValue := cnt
	prefix := fmt.Sprintf("Value=%s:%0.2f", v.Value, cnt)
	if !v.UseCount && total > 0 {
		testValue = testValue / float64(total)
		prefix = fmt.Sprintf("Value=%s:%0.2f%%", v.Value, testValue*100)
	}

	e := v.IsOkay(testValue, prefix)
	e.Description = v.GetDescription()
	return e
}
