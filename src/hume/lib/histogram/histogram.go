package histogram

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"io/ioutil"
	"strconv"
)

type histogram struct {
	Entries []histEntry `valid:"required" json:"histogram"`
}

type histEntry struct {
	Key   string  `valid:"required" json:"key"`
	Value float64 `valid:"required" json:"value"`
}

const (
	float       string = "float"
	categorical string = "categorical"
)

func GetStringMap(file string) (map[string]float64, error) {
	return GetMap(file, categorical)
}

func GetFloatMap(file string) (map[string]float64, error) {
	return GetMap(file, float)
}

func GetMap(file string, distType string) (map[string]float64, error) {
	result := make(map[string]float64)
	var err error
	//byte slice from file
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return result, err
	}
	//parse byte slice to json
	var hist histogram
	err = json.Unmarshal(data, &hist)
	if err != nil {
		return result, err
	}
	_, err = govalidator.ValidateStruct(hist)
	if err != nil {
		return result, err
	}

	for _, he := range hist.Entries {
		result[he.Key] = he.Value
	}
	if distType == float {
		for k, _ := range result {
			_, err = strconv.ParseFloat(k, 64)
			if err != nil {
				return result, err
			}
		}
	}
	return result, err
}
