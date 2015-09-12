package record

var OK = 0
var NOT_ENOUGH_VALUES = 1
var TOO_MANY_VALUES = 2

type Record struct {
	Id     string
	Number int32
	Fields *[]string
	Values []string
	Map    map[string]string
}

func (r *Record) InitMap() {
	r.Map = make(map[string]string)
	lf := len(*r.Fields)
	lv := len(r.Values)
	for i := 0; i < lf && i < lv; i++ {
		r.Map[(*r.Fields)[i]] = r.Values[i]
	}
}

func (r Record) IsValid() (bool, int) {
	if len(*r.Fields) > len(r.Values) {
		return false, NOT_ENOUGH_VALUES
	} else if len(*r.Fields) < len(r.Values) {
		return false, TOO_MANY_VALUES
	} else {
		return true, OK
	}
}

func (r Record) FieldCount() int {
	return len(*r.Fields)
}

func (r Record) ValueCount() int {
	return len(r.Values)
}
