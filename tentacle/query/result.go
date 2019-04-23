package query

import (
	"encoding/json"
)

type ResultDTO struct {
	Arguments map[string]interface{}
}

func (r ResultDTO) GetField(field string) interface{} {
	return r.Arguments[field]
}

func (r ResultDTO) HasField(field string) bool {
	_, ok := r.Arguments[field]
	return ok
}

func (r *ResultDTO) UnmarshalArguments(rawMessage []byte) error {
	if err := json.Unmarshal(rawMessage, &r.Arguments); err != nil {
		return err
	}
	return nil
}
