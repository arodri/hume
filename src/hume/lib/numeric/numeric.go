package numeric

import(
	"strconv"
)

const (
	NOT_A_FLOAT string = "NOT_A_FLOAT"
	EMPTY_STRING string = "EMPTY_STRING"
)

type FloatMapper struct {
	FloatSlice []float64
	Float2String map[float64]string
	TotalFloat float64
	TotalNonfloat float64
}

func ND_Mapper(data map[string]float64) (FloatMapper, error) {
	var FloatSlice []float64
	Float2String := make(map[float64]string)
	TotalFloat := float64(0)
	TotalNonfloat := float64(0)
	for k, v := range data {
		if k != EMPTY_STRING && k != NOT_A_FLOAT {
			f, _ := strconv.ParseFloat(k, 64)
			FloatSlice = append(FloatSlice, f)
			Float2String[f] = k
			TotalFloat += v
		} else {
			TotalNonfloat += v
		}
	}
	return FloatMapper{FloatSlice,Float2String,TotalFloat,TotalNonfloat}, nil
}
