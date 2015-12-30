package metric

import (
	"hume/lib/numeric"
)

type Mean struct {
	BaseMetric
	NumericDistribution
	result map[string]float64
}

func (m *Mean) Finalize() error {
	m.result = make(map[string]float64)

	fm, err := numeric.ND_Mapper(m.counts)
	keys := fm.FloatSlice
	f2s := fm.Float2String
	totalFloat := fm.TotalFloat

	sum := float64(0)
	for _, k := range keys {
		sum += k * m.counts[f2s[k]]
	}
	m.result["mean"] = sum / totalFloat

	return err
}

func (m *Mean) Result() MetricResult {
	return MetricResult{m.result, m.total}
}
