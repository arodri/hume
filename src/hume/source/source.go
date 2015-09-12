package source

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"hume/evaluator"
	"hume/metric"
	"hume/reader"
	"hume/record"
	"io/ioutil"
	"sync"
)

type Source struct {
	Reader       reader.Reader         `json:"reader"`
	Metrics      []metric.Metric       `json:"metrics"`
	Evaluators   []evaluator.Evaluator `json:"evaluators"`
	metricWG     sync.WaitGroup
	metricLookup map[string]int
}

func (s *Source) Init() error {
	err := s.Reader.Init()
	if err != nil {
		logrus.Fatal("Error initalizing reader", err)
	}

	s.metricLookup = make(map[string]int)
	for i, m := range s.Metrics {
		s.metricLookup[m.GetName()] = i

		m.SetInputChannel(make(chan *record.Record, 1000))
		m.SetWaitGroup(&s.metricWG)
		s.metricWG.Add(1)
		m.Init()
		m.Collect()
	}

	for _, evaluator := range s.Evaluators {
		_, found := s.metricLookup[evaluator.GetMetricName()]
		if !found {
			logrus.Fatal(fmt.Sprintf("Unknown metric: %s", evaluator.GetMetricName()))
		}
	}

	return err
}

func (s *Source) Collect() error {
	for rec := range s.Reader.ChannelIterator(1000) {
		// logrus.Debugf("%#v", rec)
		for _, metric := range s.Metrics {
			metric.GetInputChannel() <- rec
		}
	}
	for _, metric := range s.Metrics {
		close(metric.GetInputChannel())
	}
	s.metricWG.Wait()
	return nil
}

func (s *Source) getMetric(name string) metric.Metric {
	return s.Metrics[s.metricLookup[name]]
}

func (s *Source) Evaluate() (int, int) {
	err_cnt := 0
	for _, e := range s.Evaluators {
		ev := e.Evaluate(s.getMetric(e.GetMetricName()).Result())
		if ev.Ok {
			logrus.Infof("%s: OK", e.GetDescription())
		} else {
			logrus.Error(fmt.Sprintf("%s: Failed", e.GetDescription()))
			err_cnt += 1
		}
		logrus.Debug(ev.Msg)
	}
	return err_cnt, len(s.Evaluators)
}

type SourceConfig struct {
	Reader     json.RawMessage   `json:"reader"`
	Metrics    []json.RawMessage `json:"metrics"`
	Evaluators []json.RawMessage `json:"evaluators"`
}

func (sc SourceConfig) GetSource() Source {
	s := Source{
		Reader:     reader.GetReader(sc.Reader),
		Metrics:    metric.GetMetrics(sc.Metrics),
		Evaluators: evaluator.GetEvaluators(sc.Evaluators),
	}
	return s
}

func SourceConfigFromFile(file string) (SourceConfig, error) {
	sc := SourceConfig{}
	config, err := ioutil.ReadFile(file)
	if err != nil {
		return sc, err
	}
	err = json.Unmarshal(config, &sc)
	return sc, err
}
