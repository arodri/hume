package evaluator

import (
	"hume/lib/numeric"
	"hume/lib/histogram"
	"fmt"
	"sort"
	"math"
	"strconv"
)

type KS struct {
	BaseEvaluator
	Thresholder
	FileBenchmark string `json:"benchmark"`
}

func arrayMapper(array1 []float64, array2 []float64) map[float64]int {
	result := make(map[float64]int)

	for _, v := range array1 {
		result[v] = 1
	}
	for _, v := range array2 {
		result[v] += 2
	}

	return result
}

func (ks *KS) Evaluate(data map[string]float64, total int) Evaluation {

	var keys_data []float64
	map_data := make(map[float64]string)
	totalZero := float64(0)
	for k, v := range data {
		if k != numeric.EMPTY_STRING && k != numeric.NOT_A_FLOAT {
			f, _ := strconv.ParseFloat(k, 64)
			keys_data = append(keys_data, f)
			map_data[f] = k
			totalZero += v
		}
	}

	var keys_benchmark []float64
	map_benchmark:= make(map[float64]string)
	totalOne := float64(0)
	benchmark, _ := histogram.GetMap(ks.FileBenchmark)
	for k, v := range benchmark {
		if k != numeric.EMPTY_STRING && k != numeric.NOT_A_FLOAT {
			f, _ := strconv.ParseFloat(k, 64)
			keys_benchmark = append(keys_benchmark, f)
			map_benchmark[f] = k
			totalOne += v
		}
	}

	arrayMap := arrayMapper(keys_data, keys_benchmark)

	var order []float64
	for k, _ := range arrayMap {
		order = append(order, k)
	}
	sort.Sort(sort.Float64Slice(order))

	zeroSoFar := float64(0)
	oneSoFar := float64(0)
	currentKS := math.NaN()
	maxKS := float64(0)

	for _, f := range order {
		switch {
			case arrayMap[f] == 1 :
				zeroSoFar += data[map_data[f]]
			case arrayMap[f] == 2 :
				oneSoFar += benchmark[map_benchmark[f]]
			case arrayMap[f] == 3 :
				zeroSoFar += data[map_data[f]]
				oneSoFar += benchmark[map_benchmark[f]]
		}

		currentKS = (zeroSoFar/totalZero) - (oneSoFar/totalOne)

		if math.Abs(currentKS) > math.Abs(maxKS) || math.IsNaN(currentKS) {
			maxKS = currentKS
		}
	}

	testValue := maxKS * 100
	prefix := fmt.Sprintf("Value=%s:abs(%f)", "KS", testValue)
	e := ks.IsOkay(math.Abs(testValue), prefix)
	e.Description = ks.GetDescription()
	return e
}
