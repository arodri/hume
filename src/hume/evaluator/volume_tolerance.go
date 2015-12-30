package evaluator

import (
	"fmt"
)

type VolumeTolerance struct {
	BaseEvaluator
	Tolerance
}

func (v *VolumeTolerance) Evaluate(data map[string]float64, total int) Evaluation {
	testValue := float64(total)
	prefix := fmt.Sprintf("Count=%d", total)
	e := v.IsOkay(testValue, prefix)
	e.Description = v.GetDescription()
	return e
}

func (v *VolumeTolerance) Train (data map[string]float64, total int) error {
	v.SetValue(float64(total))
	return nil
}
