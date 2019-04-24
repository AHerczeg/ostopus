package query

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ostopus/tentacle/os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
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
		fetchQueryBool  bool
		executeResponse bytes.Buffer
		executeError    error
	}
	tests := []struct {
		name    string
		args    args
		mocks   mocks
		want    ResultDTO
		wantErr bool
	}{
		{
			name: "Empty query",
			args: args{name: ""},
			mocks: mocks{
				fetchQueryQuery: "",
				fetchQueryBool:  false,
			},
			want:    ResultDTO{},
			wantErr: true,
		},
		{
			name: "Normal query not saved to the query store",
			args: args{name: "Find_Containers"},
			mocks: mocks{
				fetchQueryQuery: "",
				fetchQueryBool:  false,
			},
			want:    ResultDTO{},
			wantErr: true,
		},
		{
			name: "Normal query, saved to the store",
			args: args{name: "kernel_info"},
			mocks: mocks{
				fetchQueryQuery: "SELECT * FROM kernel_info;",
				fetchQueryBool:  true,
				executeResponse: *bytes.NewBuffer([]byte(`[{"arguments":"","device":"1","path":"/System/Library/","version":"1.1.0"}]`)),
				executeError:    nil,
			},
			want: ResultDTO{
				Arguments: map[string]interface{}{
					"arguments": "",
					"device":    "1",
					"path":      "/System/Library/",
					"version":   "1.1.0",
				},
			},
			wantErr: false,
		},
		{
			name: "Faulty query, saved to the store",
			args: args{name: "kernel_info"},
			mocks: mocks{
				fetchQueryQuery: "SELECT * OH NO THIS IS NOT GOOD FROM kernel_info;",
				fetchQueryBool:  true,
				executeResponse: *bytes.NewBuffer([]byte(`[]`)),
				executeError:    nil,
			},
			want:    ResultDTO{},
			wantErr: true,
		},
		{
			name: "Normal query, saved to the store, result is not correct JSON",
			args: args{name: "kernel_info"},
			mocks: mocks{
				fetchQueryQuery: "SELECT * FROM kernel_info;",
				fetchQueryBool:  true,
				executeResponse: *bytes.NewBuffer([]byte(`Invalid:JSON:result`)),
				executeError:    nil,
			},
			want:    ResultDTO{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qs := MockQueryStore{}
			qs.On("GetQuery", mock.AnythingOfType("string")).Return(tt.mocks.fetchQueryQuery, tt.mocks.fetchQueryBool)

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

func TestStdHandler_RunCustomQuery(t *testing.T) {
	type args struct {
		query string
	}
	type mocks struct {
		fetchQueryQuery string
		fetchQueryBool  bool
		executeResponse bytes.Buffer
		executeError    error
	}
	tests := []struct {
		name    string
		args    args
		mocks   mocks
		want    ResultDTO
		wantErr bool
	}{
		{
			name: "Empty query",
			args: args{query: ""},
			mocks: mocks{
				executeResponse: bytes.Buffer{},
				executeError:    nil,
			},
			want:    ResultDTO{},
			wantErr: true,
		},
		{
			name: "Normal query",
			args: args{query: "SELECT * FROM kernel_info;"},
			mocks: mocks{
				executeResponse: *bytes.NewBuffer([]byte(`[{"arguments":"","device":"1","path":"/System/Library/","version":"1.1.0"}]`)),
				executeError:    nil,
			},
			want: ResultDTO{
				Arguments: map[string]interface{}{
					"arguments": "",
					"device":    "1",
					"path":      "/System/Library/",
					"version":   "1.1.0",
				},
			},
			wantErr: false,
		},
		{
			name: "Faulty query",
			args: args{query: "SELECT * OH NO THIS IS NOT GOOD FROM kernel_info;"},
			mocks: mocks{
				executeResponse: *bytes.NewBuffer([]byte(`[]`)),
				executeError:    nil,
			},
			want:    ResultDTO{},
			wantErr: true,
		},
		{
			name: "Normal query, result is not correct JSON",
			args: args{query: "SELECT * FROM kernel_info;"},
			mocks: mocks{
				executeResponse: *bytes.NewBuffer([]byte(`Invalid:JSON:result`)),
				executeError:    nil,
			},
			want:    ResultDTO{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qs := MockQueryStore{}

			oh := os.MockOsHandler{}
			oh.On("Execute", mock.AnythingOfType("string")).Return(tt.mocks.executeResponse, tt.mocks.executeError)

			qh := StdHandler{
				queryStore: &qs,
				osHandler:  &oh,
			}

			got, err := qh.RunCustomQuery(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("StdHandler.RunCustomQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StdHandler.RunCustomQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
