package query

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/mock"
	"ostopus/tentacle/os"
	"reflect"
	"testing"
)

func Test_cleanJSON(t *testing.T) {
	type args struct {
		out []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "Empty byte buffer",
			args: args{out: []byte{}},
			want: []byte{},
		},
		{
			name: "Clean json, no modification needed",
			args: args{out: []byte("{\"a\":\"b\"}")},
			want: []byte("{\"a\":\"b\"}"),
		},
		{
			name: "Clean json with special characters",
			args: args{out: []byte("{\"!@#$%^&*(),.\":\"b\"}")},
			want: []byte("{\"!@#$%^&*(),.\":\"b\"}"),
		},
		{
			name: "json surrounded by brackets",
			args: args{out: []byte("[{\"a\":\"b\"}]")},
			want: []byte("{\"a\":\"b\"}"),
		},
		{
			name: "Invalid json",
			args: args{out: []byte("{\"a\":\"b\":}")},
			want: []byte{},
		},
		{
			name: "Valid json with invalid special characters",
			args: args{out: []byte("{\"a//\":\"b\":}")},
			want: []byte{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buffer bytes.Buffer
			buffer.Write(tt.args.out)
			fmt.Println(json.Valid(tt.args.out))
			if got := cleanJSON(buffer); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cleanJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_RunSavedQuery(t *testing.T) {
	type args struct {
		name string
	}
	type mocks struct {
		fetchQueryQuery string
		fetchQueryError	error
		executeResponse	bytes.Buffer
		executeError	error
	}
	tests := []struct {
		name    string
		args    args
		mocks 	mocks
		want    ResultDTO
		wantErr bool
	}{
		{
			name: "",
			args: args{name: ""},
			mocks:	mocks{
				fetchQueryQuery: "",
				fetchQueryError: fmt.Errorf(""),
				executeResponse: bytes.Buffer{},
				executeError: fmt.Errorf(""),
			},
			want: ResultDTO{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qs := MockQueryStore{}
			qs.On("GetQuery", mock.AnythingOfType("string")).Return(tt.mocks.fetchQueryQuery, tt.mocks.fetchQueryError)

			oh := os.MockOsHandler{}
			oh.On("Execute", mock.AnythingOfType("string")).Return(tt.mocks.executeResponse, tt.mocks.executeError)

			qh := StdHandler{
				queryStore: &qs,
				osHandler:  &oh,
			}

			got, err := qh.RunSavedQuery(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("StdHandler.RunSavedQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StdHandler.RunSavedQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}




