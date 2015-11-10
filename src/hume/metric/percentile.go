package metric

import(
	"hume/lib/numeric"
	"strconv"
	"sort"
)

type Percentile struct {
	NumericDistribution
	result map[string]float64
	Measure int `json:"measure"`
}

//NEAREST RANK METHOD (USES ONLY MEMBERS OF POPULATION - NO INTERPOLATION)
func (p *Percentile) Finalize() error {
	p.result = make(map[string]float64)

	var keys []float64
	m := make(map[float64]string)
	for k, v := range p.counts {
		if k != numeric.EMPTY_STRING && k != numeric.NOT_A_FLOAT {
			f, _ := strconv.ParseFloat(k, 64)
			keys = append(keys, f)
			m[f] = k
		} else {
			p.total -= int(v)
		}
	}
	sort.Sort(sort.Float64Slice(keys))

	threshold := float64(p.total) * float64(p.Measure) / 100
	sum := float64(0)
	//keep index of key coorsponding to last sum
	var i int
	var k float64
	for i, k = range keys {
		sum += float64(p.counts[m[k]])
		if sum >= threshold {
			break
		}
	}
	p.result["percentile"] = keys[i]

	return nil
}

func (p *Percentile) Result() MetricResult{
	return MetricResult{p.result, p.total}
}
