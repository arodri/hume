package evaluator

import (
	"fmt"
	"strings"
	"math"
)

type Thresholder struct {
	Threshold      float64 `json:"threshold"`
	ErrorOn        string  `json:"error_on"`
	errorOnLess    bool    `json:"error_less"`
	errorOnEqual   bool    `json:"error_equal"`
	errorOnGreater bool    `json:"error_greater"`
}

func (t *Thresholder) IsOkay(testValue float64, msg_prefix string) Evaluation {
	e := Evaluation{
		Msg: "",
		Ok:  true,
	}

	errorOnLess := strings.Contains(t.ErrorOn, "<")
	errorOnEqual := strings.Contains(t.ErrorOn, "=")
	errorOnGreater := strings.Contains(t.ErrorOn, ">")

	if errorOnEqual && testValue == t.Threshold {
		e.Ok = false
	} else if errorOnGreater && testValue > t.Threshold {
		e.Ok = false
	} else if errorOnLess && testValue < t.Threshold {
		e.Ok = false
	} else if math.IsNaN(testValue) { 
		e.Ok = false
	}

	p1, p2, p3 := "", "", ""
	if errorOnLess {
		p1 = "<"
	}
	if errorOnEqual {
		p2 = "="
	}
	if errorOnGreater {
		p3 = ">"
	}
	e.Msg = fmt.Sprintf(
		"%s ErrorOn: M %s%s%s T, Threshold: %0.2f, Measured: %0.2f",
		msg_prefix, p1, p2, p3, t.Threshold, testValue,
	)

	return e
}
