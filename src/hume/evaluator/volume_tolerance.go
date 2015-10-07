package evaluator

import (
	"fmt"
)

type VolumeTolerance struct {
	BaseEvaluator
	Tolerance
}

func (v VolumeTolerance) Evaluate(data map[string]int, total int) Evaluation {
	testValue := float64(total)
	prefix := fmt.Sprintf("Count=%d", total)
	e := v.IsOkay(testValue, prefix)
	e.Description = v.GetDescription()
	return e
}
