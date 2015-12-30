package histogram

import (
	"github.com/asaskevich/govalidator"
	"strconv"
)

type HistEntry struct {
	Key   string  `valid:"required" json:"key"`
	Value float64 `valid:"required" json:"value"`
}

const (
	continuous  string = "continuous"
	categorical string = "categorical"
)

func GetStringMap(histogram []HistEntry) (map[string]float64, error) {
	return GetMap(histogram, categorical)
}

func GetFloatMap(histogram []HistEntry) (map[string]float64, error) {
	return GetMap(histogram, continuous)
}

func GetMap(histogram []HistEntry, distType string) (map[string]float64, error) {
	result := make(map[string]float64)
	var err error

	for _, entry := range histogram {
		_, err = govalidator.ValidateStruct(entry)
		if err != nil {
			return result, err
		}
		result[entry.Key] = entry.Value
	}
	if distType == continuous {
		for k, _ := range result {
			_, err = strconv.ParseFloat(k, 64)
			if err != nil {
				return result, err
			}
		}
	}
	return result, err
}

func MapToHist(data map[string]float64) ([]HistEntry, error) {
	var histogram []HistEntry
	var entry HistEntry
	for k, v := range data {
		entry.Key = k
		entry.Value = v
		histogram = append(histogram, entry)
	}
	return histogram, nil
}
