package source

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"hume/metric"
	"hume/reader"
	"hume/record"
	"io/ioutil"
	"sync"
)

type Source struct {
	Reader   reader.Reader   `json:"reader"`
	Metrics  []metric.Metric `json:"metrics"`
	metricWG sync.WaitGroup
}

func (s *Source) Init() error {
	err := s.Reader.Init()
	if err != nil {
		logrus.Fatal("Error initalizing reader", err)
	}

	for _, m := range s.Metrics {
		logrus.Debug(fmt.Sprintf("Initializing Metric: %s", m.GetName()))
		m.SetInputChannel(make(chan *record.Record, 1000))
		m.SetWaitGroup(&s.metricWG)

		err = m.Init()
		if err != nil {
			logrus.Fatal(fmt.Sprintf("Error intializing %s", m.GetName()), err)
		}

		metric.Collect(m)
	}

	return err
}

func (s *Source) Collect() error {
	var err error
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
	for _, metric := range s.Metrics {
		err = metric.Finalize()
		if err != nil {
			return err
		}
	}
	return err
}

func (s *Source) Evaluate() (int, int) {
	err_cnt := 0
	total := 0

	for _, m := range s.Metrics {
		for _, e := range metric.Evaluate(m) {
			if e.Ok {
				logrus.Infof("%s: OK", e.Description)
			} else {
				logrus.Error(fmt.Sprintf("%s: Failed", e.Description))
				err_cnt += 1
			}
			total += 1
			logrus.Debug(e.Msg)
		}
	}
	return err_cnt, total
}

func (s *Source) Train() (int, int) {
	err_cnt := 0
	total := 0

	for _, m := range s.Metrics {
		logrus.Debugf("Training: %s", m.GetName())
		this_cnt, this_total := metric.Train(m)
		err_cnt += this_cnt
		total += this_total
	}
	return err_cnt, total
}

type SourceConfig struct {
	Reader  json.RawMessage   `json:"reader"`
	Metrics []json.RawMessage `json:"metrics"`
}

func (sc SourceConfig) GetSource() Source {
	s := Source{
		Reader:  reader.GetReader(sc.Reader),
		Metrics: metric.GetMetrics(sc.Metrics),
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
