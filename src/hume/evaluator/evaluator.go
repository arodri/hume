package evaluator

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"hume/metric"
)

func GetEvaluators(configs []json.RawMessage) []Evaluator {
	evaluators := []Evaluator{}
	for _, config := range configs {
		evaluators = append(evaluators, GetEvaluator(config))
	}
	return evaluators
}

func GetEvaluator(config []byte) Evaluator {
	var err error

	be := BaseEvaluator{}
	err = json.Unmarshal(config, &be)
	if err != nil {
		log.Fatal("Invalid evaluator config, need evaluator, metric_name, and silence defined", err)
	}

	var e Evaluator
	switch be.Evaluator {
	default:
		log.Fatal(fmt.Sprintf("Unrecognized evaluator: %#v", be))
	case "ValueThreshold":
		v := ValueThreshold{}
		err = json.Unmarshal(config, &v)
		e = &v
	case "VolumeTolerance":
		v := VolumeTolerance{}
		err = json.Unmarshal(config, &v)
		e = &v
	}

	if err != nil {
		log.Fatal(fmt.Sprintf("Invalid configuration for %#v", be))
	}

	return e
}

type BaseEvaluator struct {
	Evaluator   string `json:"evaluator"`
	MetricName  string `json:"metric_name"`
	Silence     bool   `json:"silence"`
	Description string `json:"description"`
}

func (be *BaseEvaluator) GetName() string {
	return be.Evaluator
}

func (be *BaseEvaluator) ShouldAlert() bool {
	return be.Silence
}

func (be *BaseEvaluator) GetMetricName() string {
	return be.MetricName
}

func (be *BaseEvaluator) GetDescription() string {
	return be.Description
}

type Evaluator interface {
	Evaluate(test metric.MetricResult) Evaluation
	GetName() string
	ShouldAlert() bool
	GetMetricName() string
	GetDescription() string
}

type Evaluation struct {
	Msg string
	Ok  bool
}
