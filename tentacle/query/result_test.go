package query

import (
	"encoding/json"
	"reflect"
	"testing"
	"testing/quick"
)

func TestResultDTO_GetField(t *testing.T) {
	f := func(dto ResultDTO) bool {
		for k := range dto.Arguments {
			return reflect.DeepEqual(dto.GetField(k), dto.Arguments[k])
		}
		return true
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestResultDTO_HasField(t *testing.T) {
	f := func(dto ResultDTO) bool {
		for k := range dto.Arguments {
			return dto.HasField(k)
		}
		return true
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestResultDTO_UnmarshalArguments(t *testing.T) {
	f := func(result map[string]string) bool {
		m, err := json.Marshal(result)
		if err != nil {
			// this err is not relevant for the test
			return true
		}
		var dto ResultDTO
		err = dto.UnmarshalArguments(m)
		return err == nil
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
