package evaluator

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
)

func GetEvaluators(configs []json.RawMessage) []Evaluator {
	evaluators := []Evaluator{}
	for _, config := range configs {
		e := GetEvaluator(config)
		log.Debug(fmt.Sprintf("Created Evaluator: %s", e.GetDescription()))
		evaluators = append(evaluators, GetEvaluator(config))
	}
	return evaluators
}

func GetEvaluator(config []byte) Evaluator {
	var err error

	be := BaseEvaluator{}
	err = json.Unmarshal(config, &be)
	if err != nil {
		log.Error(err)
		log.Fatal("Invalid evaluator config, need evaluator, and silence defined", err)
	}

	var e Evaluator
	switch be.Evaluator {
	default:
		log.Fatal(fmt.Sprintf("Unrecognized evaluator: %#v", be))
	case "ValueThreshold":
		v := ValueThreshold{}
		err = json.Unmarshal(config, &v)
		e = &v
	case "ValueTolerance":
		v := ValueThreshold{}
		err = json.Unmarshal(config, &v)
		e = &v
	case "VolumeTolerance":
		v := VolumeTolerance{}
		err = json.Unmarshal(config, &v)
		e = &v
	case "NearestValueThreshold":
		v := NearestValueThreshold{}
		err = json.Unmarshal(config, &v)
		e = &v
	case "NearestValueTolerance":
		v := NearestValueTolerance{}
		err = json.Unmarshal(config, &v)
		e = &v
	case "KS":
		v := KS{}
		err = json.Unmarshal(config, &v)
		e = &v
	}

	if err != nil {
		log.Error(err)
		log.Fatal(fmt.Sprintf("Invalid configuration for %#v", be))
	}

	return e
}

type BaseEvaluator struct {
	Evaluator   string `json:"evaluator"`
	Silence     bool   `json:"silence"`
	Description string `json:"description"`
}

func (be *BaseEvaluator) GetType() string {
	return be.Evaluator
}

func (be *BaseEvaluator) ShouldAlert() bool {
	return be.Silence
}

func (be *BaseEvaluator) GetDescription() string {
	return be.Description
}

type Evaluator interface {
	Evaluate(data map[string]float64, total int) Evaluation
	Train(data map[string]float64, total int) error
	GetType() string
	ShouldAlert() bool
	GetDescription() string
}

type Evaluation struct {
	Description string
	Msg         string
	Ok          bool
}
