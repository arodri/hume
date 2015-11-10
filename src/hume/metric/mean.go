package metric

import (
	"hume/lib/numeric"
	"strconv"
)

type Mean struct {
	NumericDistribution
	result map[string]float64
}

func (m *Mean) Finalize() error {
	m.result = make(map[string]float64)
	sum := float64(0)
	for k, v := range m.counts {
		if k != numeric.EMPTY_STRING && k != numeric.NOT_A_FLOAT {
			f, _ := strconv.ParseFloat(k, 64)
			sum += f * v
		} else {
			m.total -= int(v)
		}
	}
	m.result["mean"] = sum / float64(m.total)

	return nil
}

func (m *Mean) Result() MetricResult {
	return MetricResult{m.result, m.total}
}
