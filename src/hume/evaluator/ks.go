package evaluator

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"hume/lib/histogram"
	"hume/lib/numeric"
	"math"
	"sort"
)

type KS struct {
	BaseEvaluator
	Thresholder
	BenchmarkJSON []histogram.HistEntry `json:"benchmark"`
}

func arrayMapper(array1 []float64, array2 []float64) (map[float64][2]bool, []float64) {
	mapper := make(map[float64][2]bool)
	var keys []float64

	for _, v := range array1 {
		mapper[v] = [2]bool{true, false}
		keys = append(keys, v)
	}
	for _, v := range array2 {
		b := mapper[v]
		b[1] = true
		mapper[v] = b
		if b[0] == false {
			keys = append(keys, v)
		}
	}

	return mapper, keys
}

func (ks *KS) compute(data map[string]float64, total int) (float64, error, string) {
	var err error
	var msg string
	testValue := float64(0)

	fm_data, _ := numeric.ND_Mapper(data)
	keys_data := fm_data.FloatSlice
	map_data := fm_data.Float2String
	totalZero := fm_data.TotalFloat

	benchmark, err := histogram.GetFloatMap(ks.BenchmarkJSON)
	if err != nil {
		msg = fmt.Sprintf("Error in benchmark float map: %s", err)
		logrus.Error(msg)
		return testValue, err, msg
	}
	fm_benchmark, err := numeric.ND_Mapper(benchmark)
	if err != nil {
		msg = fmt.Sprintf("Error in mapping benchmark to NumericDistribution: %s", err)
		logrus.Error(msg)
		return testValue, err, msg
	}
	keys_benchmark := fm_benchmark.FloatSlice
	map_benchmark := fm_benchmark.Float2String
	totalOne := fm_benchmark.TotalFloat

	key_map, keys_all := arrayMapper(keys_data, keys_benchmark)
	sort.Sort(sort.Float64Slice(keys_all))

	zeroSoFar := float64(0)
	oneSoFar := float64(0)
	currentKS := math.NaN()
	maxKS := float64(0)

	for _, f := range keys_all {
		if key_map[f][0] {
			zeroSoFar += data[map_data[f]]
		}
		if key_map[f][1] {
			oneSoFar += benchmark[map_benchmark[f]]
		}

		currentKS = (zeroSoFar / totalZero) - (oneSoFar / totalOne)

		if math.Abs(currentKS) > math.Abs(maxKS) || math.IsNaN(currentKS) {
			maxKS = currentKS
		}
	}

	testValue = maxKS * 100
	return testValue, err, msg
}

func (ks *KS) Evaluate(data map[string]float64, total int) Evaluation {
	e := Evaluation{ks.GetDescription(), "", false}
	testValue, err, msg := ks.compute(data, total)
	if err != nil {
		e.Msg = msg
		return e
	}

	prefix := fmt.Sprintf("Value=%s:abs(%f)", "KS", testValue)
	e = ks.IsOkay(math.Abs(testValue), prefix)
	e.Description = ks.GetDescription()
	return e
}

func (ks *KS) SetBenchmark(histogram []histogram.HistEntry) {
	ks.BenchmarkJSON = histogram
}

func (ks *KS) Train(data map[string]float64, total int) error {
	hist, err := histogram.MapToHist(data)
	ks.SetBenchmark(hist)
	return err
}
