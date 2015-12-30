package metric

import (
	"hume/lib/numeric"
	"sort"
	"strconv"
)

type Percentile struct {
	BaseMetric
	NumericDistribution
	result map[string]float64
}

//NEAREST RANK METHOD (USES ONLY MEMBERS OF POPULATION - NO INTERPOLATION)
func (p *Percentile) Finalize() error {
	p.result = make(map[string]float64)

	fm, _ := numeric.ND_Mapper(p.counts)
	keys := fm.FloatSlice
	m := fm.Float2String
	totalFloat := fm.TotalFloat
	sort.Sort(sort.Float64Slice(keys))

	for c := float64(0); c < 100; c++ {
		threshold := totalFloat * c / 100
		sum := float64(0)
		//keep index of key coorsponding to last sum
		var i int
		var k float64
		for i, k = range keys {
			sum += p.counts[m[k]]
			if sum >= threshold {
				break
			}
		}
		p.result[strconv.FormatFloat(c, 'f', -1, 64)] = keys[i]
	}

	return nil
}

func (p *Percentile) Result() MetricResult {
	return MetricResult{p.result, p.total}
}
