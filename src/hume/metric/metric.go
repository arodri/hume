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

	m := BaseMetric{}
	err = json.Unmarshal(config, &m)
	if err != nil {
		logrus.Error(err)
		logrus.Fatal(fmt.Sprintf("Missing metric name and/or type"))
	}

	switch m.Type {
	default:
		logrus.Fatal(fmt.Sprintf("Unrecognized metric type: %#v", m))
	case "FieldCount":
		ft := FieldCount{}
		err = json.Unmarshal(config, &ft)
		m.SetMetricImpl(&ft)
	case "DateFormat":
		ft := DateFormat{}
		err = json.Unmarshal(config, &ft)
		m.SetMetricImpl(&ft)
	case "Populated":
		ft := Populated{}
		err = json.Unmarshal(config, &ft)
		m.SetMetricImpl(&ft)
	case "FutureDate":
		ft := FutureDate{}
		err = json.Unmarshal(config, &ft)
		m.SetMetricImpl(&ft)
	case "RegexTest":
		ft := RegexTest{}
		err = json.Unmarshal(config, &ft)
		m.SetMetricImpl(&ft)
	case "NumRecords":
		ft := NumRecords{}
		err = json.Unmarshal(config, &ft)
		m.SetMetricImpl(&ft)
	case "ValidValues":
		ft := ValidValues{}
		err = json.Unmarshal(config, &ft)
		m.SetMetricImpl(&ft)
	case "NominalDistribution":
		ft := NominalDistribution{}
		err = json.Unmarshal(config, &ft)
		m.SetMetricImpl(&ft)
	case "NumericDistribution":
		ft := NumericDistribution{}
		err = json.Unmarshal(config, &ft)
		m.SetMetricImpl(&ft)
	case "Mean":
		ft := Mean{}
		err = json.Unmarshal(config, &ft)
		m.SetMetricImpl(&ft)
	case "Percentile":
		ft := Percentile{}
		err = json.Unmarshal(config, &ft)
		m.SetMetricImpl(&ft)

	}
	if err != nil {
		logrus.Error(err)
		logrus.Fatal(fmt.Sprintf("Invalid configuration for %#v", m))
	}

	return &m
}

type BaseMetric struct {
	Name             string            `json:"name"`
	Type             string            `json:"type"`
	EvaluatorConfigs []json.RawMessage `json:"evaluators"`
	evaluators       []evaluator.Evaluator
	inputChannel     chan *record.Record
	waitGroup        *sync.WaitGroup
	metricImpl       MetricImpl
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

func (bm *BaseMetric) SetMetricImpl(m MetricImpl) {
	bm.metricImpl = m
}

func (bm *BaseMetric) SetEvaluators(e []evaluator.Evaluator) {
	bm.evaluators = e
}

func (bm *BaseMetric) Collect() {
	go func() {
		for rec := range bm.inputChannel {
			bm.mutex.Lock()
			bm.metricImpl.Process(rec)
			bm.mutex.Unlock()
		}
		bm.waitGroup.Done()
	}()
}

func (bm *BaseMetric) Evaluate() []evaluator.Evaluation {
	evaluations := []evaluator.Evaluation{}

	mr := bm.Result()

	for _, e := range bm.evaluators {
		logrus.Debug(fmt.Sprintf("Evaluating: %s", e.GetDescription()))
		evaluations = append(evaluations, e.Evaluate(mr.Data, mr.Total))
	}
	return evaluations
}

func (bm *BaseMetric) Init() error {
	return bm.metricImpl.Init()
}

func (bm *BaseMetric) Finalize() error {
	return bm.metricImpl.Finalize()
}

func (bm *BaseMetric) Result() MetricResult {
	return bm.metricImpl.Result()
}

func (bm *BaseMetric) Process(rec *record.Record) {
	bm.metricImpl.Process(rec)
}

func (bm *BaseMetric) GetEvaluators() []evaluator.Evaluator {
	return bm.evaluators
}

func (bm *BaseMetric) GetEvaluatorConfigs() []json.RawMessage {
	return bm.EvaluatorConfigs
}

type MetricImpl interface {
	Init() error
	Process(rec *record.Record)
	Finalize() error
	Result() MetricResult
}

type Metric interface {
	MetricImpl
	GetName() string
	GetType() string
	SetWaitGroup(wg *sync.WaitGroup)
	GetWaitGroup() *sync.WaitGroup
	SetInputChannel(inputChannel chan *record.Record)
	GetInputChannel() chan *record.Record
	SetMetricImpl(m MetricImpl)
	Collect()
	Evaluate() []evaluator.Evaluation
	GetEvaluators() []evaluator.Evaluator
	SetEvaluators([]evaluator.Evaluator)
	GetEvaluatorConfigs() []json.RawMessage
}

type MetricResult struct {
	Data  map[string]float64
	Total int
}
