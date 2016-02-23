package models

import "encoding/json"

type Options struct {
	ClientId string `json:"clientId"`
}

func (m Options) ToJson() (JSON, error) {
	var model JSON

	if err := model.FromObject(m); err != nil {
		return model, err
	}

	return model, nil
}

func (m *Options) FromString(s string) error {
	if err := json.Unmarshal([]byte(s), m); err != nil {
		return err
	}

	return nil
}
