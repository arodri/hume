package evaluator

import (
	"fmt"
)

type ValueThreshold struct {
	BaseEvaluator
	Thresholder
	Value    string `json:"value"`
	UseCount bool   `json:"use_count"`
}

func (v ValueThreshold) Evaluate(data map[string]int, total int) Evaluation {
	cnt, ok := data[v.Value]
	if !ok {
		cnt = 0
	}

	testValue := float64(cnt)
	prefix := fmt.Sprintf("Value=%s:%d", v.Value, cnt)
	if !v.UseCount && total > 0 {
		testValue = testValue / float64(total)
		prefix = fmt.Sprintf("Value=%s:%0.2f%%", v.Value, testValue*100)
	}

	e := v.IsOkay(testValue, prefix)
	e.Description = v.GetDescription()
	return e
}
