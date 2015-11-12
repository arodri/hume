package evaluator

import (
	"fmt"
	"math"
)

type Tolerance struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

func (t *Tolerance) IsOkay(testValue float64, msg_prefix string) Evaluation {
	e := Evaluation{
		Msg: "",
		Ok:  true,
	}

	if testValue < t.Min || t.Max < testValue {
		e.Ok = false
	} else if math.IsNaN(testValue) {
		e.Ok = false
	}

	e.Msg = fmt.Sprintf(
		"%s Tolerance: %0.2f <= M <= %0.2f, Measured: %0.2f",
		msg_prefix, t.Min, t.Max, testValue,
	)

	return e
}
