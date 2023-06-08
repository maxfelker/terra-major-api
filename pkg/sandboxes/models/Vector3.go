package models

import (
	"encoding/json"
	"errors"
)

type Vector3 struct {
	X *float32 `json:"x"`
	Y *float32 `json:"y"`
	Z *float32 `json:"z"`
}

func (v *Vector3) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []bytes")
	}
	err := json.Unmarshal(asBytes, v)
	if err != nil {
		return err
	}
	return nil
}
