package evaluator

import (
	"fmt"
	"hume/metric"
)

type ValueThreshold struct {
	BaseEvaluator
	Thresholder
	Value    string `json:"value"`
	UseCount bool   `json:"use_count"`
}

func (v ValueThreshold) Evaluate(test metric.MetricResult) Evaluation {
	cnt, ok := test.Data[v.Value]
	if !ok {
		cnt = 0
	}
	testValue := float64(cnt)
	prefix := fmt.Sprintf("Value=%s:%d", v.Value, cnt)
	if !v.UseCount && test.Total > 0 {
		testValue = testValue / float64(test.Total)
		prefix = fmt.Sprintf("Value=%s:%0.2f%%", v.Value, testValue*100)
	}

	return v.IsOkay(testValue, prefix)
}
