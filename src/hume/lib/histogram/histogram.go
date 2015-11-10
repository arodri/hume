package histogram

import(
	"io/ioutil"
	"encoding/json"
	"fmt"
)

type HistEntry struct {
	Key string `json:"key"`
	Value float64 `json:"value"`
}

func GetMap(file string) (map[string]float64, error) {
	result := make(map[string]float64)
	var err error
	//byte slice from file
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return result, err
	}
	//parse byte slice to json
	var entries []HistEntry
	err = json.Unmarshal(data, &entries)
	if err != nil {
		fmt.Println(err)
		return result, err
	}
	for _, he := range entries {
		result[he.Key] = he.Value
	}
	return result, err
}
