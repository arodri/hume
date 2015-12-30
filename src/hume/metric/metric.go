package metric

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"hume/evaluator"
	"hume/record"
	"sync"
)

func GetMetrics(configs []json.RawMessage) []Metric {
	metrics := []Metric{}
	for _, config := range configs {
		m := GetMetric(config)
		logrus.Debug(fmt.Sprintf("Created Metric: %s", m.GetName()))
		m.SetEvaluators(evaluator.GetEvaluators(m.GetEvaluatorConfigs()))
		metrics = append(metrics, m)
	}
	return metrics
}

func GetMetric(config []byte) Metric {
	var err error

	bm := BaseMetric{}
	err = json.Unmarshal(config, &bm)
	if err != nil {
		logrus.Error(err)
		logrus.Fatal(fmt.Sprintf("Missing metric name and/or type"))
	}

	var m Metric
	switch bm.Type {
	default:
		logrus.Fatal(fmt.Sprintf("Unrecognized metric type: %#v", bm))
	case "FieldCount":
		ft := FieldCount{}
		err = json.Unmarshal(config, &ft)
		m = &ft
	case "DateFormat":
		ft := DateFormat{}
		err = json.Unmarshal(config, &ft)
		m = &ft
	case "Populated":
		ft := Populated{}
		err = json.Unmarshal(config, &ft)
		m = &ft
	case "FutureDate":
		ft := FutureDate{}
		err = json.Unmarshal(config, &ft)
		m = &ft
	case "RegexTest":
		ft := RegexTest{}
		err = json.Unmarshal(config, &ft)
		m = &ft
	case "NumRecords":
		ft := NumRecords{}
		err = json.Unmarshal(config, &ft)
		m = &ft
	case "ValidValues":
		ft := ValidValues{}
		err = json.Unmarshal(config, &ft)
		m = &ft
	case "NominalDistribution":
		ft := NominalDistribution{}
		err = json.Unmarshal(config, &ft)
		m = &ft
	case "NumericDistribution":
		ft := NumericDistribution{}
		err = json.Unmarshal(config, &ft)
		m = &ft
	case "Mean":
		ft := Mean{}
		err = json.Unmarshal(config, &ft)
		m = &ft
	case "Percentile":
		ft := Percentile{}
		err = json.Unmarshal(config, &ft)
		m = &ft
	case "CDF":
		ft := CDF{}
		err = json.Unmarshal(config, &ft)
		m = &ft
	}

	if err != nil {
		logrus.Error(err)
		logrus.Fatal(fmt.Sprintf("Invalid configuration for %#v", m))
	}

	return m
}

type BaseMetric struct {
	Name             string            `json:"name"`
	Type             string            `json:"type"`
	EvaluatorConfigs []json.RawMessage `json:"evaluators"`
	evaluators       []evaluator.Evaluator
	inputChannel     chan *record.Record
	waitGroup        *sync.WaitGroup
	mutex            sync.Mutex
}

func (bm *BaseMetric) GetName() string {
	return bm.Name
}

func (bm *BaseMetric) GetType() string {
	return bm.Type
}

func (bm *BaseMetric) SetInputChannel(inputChannel chan *record.Record) {
	bm.inputChannel = inputChannel
}

func (bm *BaseMetric) GetInputChannel() chan *record.Record {
	return bm.inputChannel
}

func (bm *BaseMetric) SetWaitGroup(wg *sync.WaitGroup) {
	bm.waitGroup = wg
}

func (bm *BaseMetric) GetWaitGroup() *sync.WaitGroup {
	return bm.waitGroup
}

func (bm *BaseMetric) SetEvaluators(e []evaluator.Evaluator) {
	bm.evaluators = e
}

func (bm *BaseMetric) Lock() {
	bm.mutex.Lock()
}

func (bm *BaseMetric) Unlock() {
	bm.mutex.Unlock()
}

func Collect(m Metric) {
	go func() {
		m.GetWaitGroup().Add(1)
		for rec := range m.GetInputChannel() {
			m.Lock()
			m.Process(rec)
			m.Unlock()
		}
		m.GetWaitGroup().Done()
	}()
}

func Evaluate(m Metric) []evaluator.Evaluation {
	evaluations := []evaluator.Evaluation{}

	mr := m.Result()

	for _, e := range m.GetEvaluators() {
		logrus.Debug(fmt.Sprintf("Evaluating: %s", e.GetDescription()))
		evaluations = append(evaluations, e.Evaluate(mr.Data, mr.Total))
	}
	return evaluations
}

func Train(m Metric) (int, int) {
	mr := m.Result()
	num_errors := 0
	evals := m.GetEvaluators()
	evals_config := m.GetEvaluatorConfigs()
	new_eval_config := make([]json.RawMessage, len(evals))
	for i, e := range evals {
		logrus.Debugf("Training Evaluator: %s", e.GetDescription())
		err := e.Train(mr.Data, mr.Total)
		if err != nil {
			num_errors += 1
			continue
		}
		j, err := json.Marshal(e)
		if err != nil {
			num_errors += 1
			new_eval_config[i] = evals_config[i]
			continue
		}
		raw := json.RawMessage(j)
		new_eval_config[i] = raw
	}

	m.SetEvaluatorConfigs(new_eval_config)
	return num_errors, len(evals)

}

func (bm *BaseMetric) GetEvaluators() []evaluator.Evaluator {
	return bm.evaluators
}

func (bm *BaseMetric) GetEvaluatorConfigs() []json.RawMessage {
	return bm.EvaluatorConfigs
}

func (bm *BaseMetric) SetEvaluatorConfigs(j []json.RawMessage) {
	bm.EvaluatorConfigs = j
}

type Metric interface {
	Init() error
	Process(rec *record.Record)
	Finalize() error
	Result() MetricResult
	// these methods are satisfied by BaseMetric
	GetName() string
	GetType() string
	SetWaitGroup(wg *sync.WaitGroup)
	GetWaitGroup() *sync.WaitGroup
	SetInputChannel(inputChannel chan *record.Record)
	GetInputChannel() chan *record.Record
	GetEvaluators() []evaluator.Evaluator
	SetEvaluators([]evaluator.Evaluator)
	GetEvaluatorConfigs() []json.RawMessage
	SetEvaluatorConfigs([]json.RawMessage)
	Lock()
	Unlock()
}

type MetricResult struct {
	Data  map[string]float64
	Total int
}
