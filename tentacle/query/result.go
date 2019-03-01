package query

import (
	"encoding/json"
)

type ResultDTO struct {
	QueryID   string
	DeviceID  string
	arguments map[string]interface{}
}

func (r ResultDTO) GetField(field string) interface{} {
	return r.arguments[field]
}

func (r ResultDTO) HasField(field string) bool {
	_, ok := r.arguments[field]
	return ok
}

func (r *ResultDTO) UnmarshalArguments(rawMessage []byte) error {
	if err := json.Unmarshal(rawMessage, &r.arguments); err != nil {
		return err
	}
	return nil
}
