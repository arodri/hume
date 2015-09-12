package evaluator

import (
	"fmt"
	"hume/metric"
)

type VolumeTolerance struct {
	BaseEvaluator
	Tolerance
}

func (v VolumeTolerance) Evaluate(test metric.MetricResult) Evaluation {
	testValue := float64(test.Total)
	prefix := fmt.Sprintf("Count=%d", test.Total)
	return v.IsOkay(testValue, prefix)
}
