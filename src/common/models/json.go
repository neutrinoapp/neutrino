package models

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/neutrinoapp/neutrino/src/common/log"
)

type JSON map[string]interface{}

func (j JSON) ForEach(f func(key string, value interface{})) {
	for k := range j {
		f(k, j[k])
	}
}

func (j JSON) String() (string, error) {
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

func (j *JSON) FromObject(o interface{}) error {
	data, err := json.Marshal(o)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, j); err != nil {
		return err
	}

	return nil
}

func (j *JSON) FromResponse(res *http.Response) error {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if string(body) == "" {
		return nil
	}

	err = json.Unmarshal(body, j)
	if err != nil {
		return errors.New(err.Error() + "; Body: " + string(body))
	}

	return nil
}

func Convert(input interface{}, target interface{}) error {
	if input == nil {
		return nil
	}

	b, err := json.Marshal(input)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, target)
}

func String(input interface{}) string {
	b, err := json.Marshal(input)
	if err != nil {
		log.Error(err)
		return "{}"
	}

	return string(b)
}
