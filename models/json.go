package models

import (
	"encoding/json"
	"net/http"
	"io/ioutil"
	"errors"
)

type JSON map[string]interface{}

func (j JSON) ForEach(f func(key string, value interface{})) {
	for k := range j {
		f(k, j[k])
	}
}

func (j *JSON) String() (string, error) {
	b, err := json.Marshal(j)
	return string(b), err
}

func (j *JSON) FromString(str []byte) error {
	return json.Unmarshal(str, j)
}

func (j JSON) FromMap(m map[string]interface{}) JSON {
	for k := range m {
		j[k] = m[k]
	}

	return j
}

func (j *JSON) FromResponse(res *http.Response) error {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if (string(body) == "") {
		return nil
	}

	err = json.Unmarshal(body, j)
	if err != nil {
		return errors.New(err.Error() + "; Body: " + string(body))
	}

	return nil
}
