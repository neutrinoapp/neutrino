package api

import (
	"bytes"
	"fmt"
)

type JSON map[string]interface{}

func (j JSON) ForEach(f func(key string, value interface{})) {
	for k := range j {
		f(k, j[k])
	}
}

func (j JSON) String() string {
	var b bytes.Buffer

	j.ForEach(func(k string, v interface{}) {
		b.WriteString(k + ":" + fmt.Sprintf("%v", v) + "\r\n")
	})

	return b.String()
}

func (j JSON) FromMap(m map[string]interface{}) JSON {
	for k := range m {
		j[k] = m[k]
	}

	return j
}
