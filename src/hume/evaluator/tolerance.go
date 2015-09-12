package evaluator

import (
	"fmt"
)

type Tolerance struct {
	Target       float64 `json:"target"`
	MinTolerance float64 `json:"min"`
	MaxTolerance float64 `json:"max"`
}

func (t *Tolerance) IsOkay(testValue float64, msg_prefix string) Evaluation {
	e := Evaluation{"", true}

	min := t.Target * t.MinTolerance
	max := t.Target * t.MaxTolerance

	if testValue < min || max < testValue {
		e.Ok = false
	}

	e.Msg = fmt.Sprintf(
		"%s Tolerance: %0.2f <= M <= %0.2f, Measured: %0.2f",
		msg_prefix, min, max, testValue,
	)

	return e
}
