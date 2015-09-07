package api

type JSON map[string]interface{}

func (j JSON) FromMap(m map[string]interface{}) JSON {
	for k := range m {
		j[k] = m[k]
	}

	return j
}