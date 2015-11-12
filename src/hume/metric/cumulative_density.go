package metric

import (
	"hume/lib/numeric"
	"sort"
)

type CDF struct {
	NumericDistribution
	result map[string]float64
}

func (cdf *CDF) Finalize() error {
	cdf.result = make(map[string]float64)

	fm, _ := numeric.ND_Mapper(cdf.counts)
	keys := fm.FloatSlice
	m := fm.Float2String
	totalFloat := fm.TotalFloat
	sort.Sort(sort.Float64Slice(keys))

	sum := float64(0)
	for _, f := range keys {
		sum += cdf.counts[m[f]]
		cdf.result[m[f]] = sum / totalFloat
	}

	return nil
}

func (cdf *CDF) Result() MetricResult {
	return MetricResult{cdf.result, cdf.total}
}
