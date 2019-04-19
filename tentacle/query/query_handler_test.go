package query

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ostopus/tentacle/os"
	"reflect"
	"testing"
)

func TestNewQueryHandler(t *testing.T) {
	type args struct {
		store QueryStore
		os    os.Handler
	}
	tests := []struct {
		name string
		args args
		want Handler
	}{
		{
			name: "Empty handlers",
			args: args{
				store: &localQueryStore{},
				os:    os.NewOSHandler(),
			},
			want: Handler{queryStore: &localQueryStore{}, osHandler: os.NewOSHandler()},
		},
		{
			name: "Empty handlers",
			args: args{
				store: &localQueryStore{},
				os:    os.NewOSHandler(),
			},
			want: Handler{queryStore: &localQueryStore{}, osHandler: os.NewOSHandler()},
		},
		{
			name: "Empty handlers",
			args: args{
				store: &localQueryStore{},
				os:    os.NewOSHandler(),
			},
			want: Handler{queryStore: &localQueryStore{}, osHandler: os.NewOSHandler()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitQueryHandler(tt.args.store, tt.args.os); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitQueryHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryHandler_RunCustomQuery(t *testing.T) {
	type fields struct {
		queryStore QueryStore
		osHandler  os.Handler
	}
	type args struct {
		query string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ResultDTO
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qh := Handler{
				queryStore: tt.fields.queryStore,
				osHandler:  tt.fields.osHandler,
			}
			got, err := qh.RunCustomQuery(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryHandler.RunCustomQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryHandler.RunCustomQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func TestQueryHandler_RunSavedQuery(t *testing.T) {
	type fields struct {
		queryStore QueryStore
		osHandler  os.Handler
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ResultDTO
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qh := Handler{
				queryStore: tt.fields.queryStore,
				osHandler:  tt.fields.osHandler,
			}
			got, err := qh.RunSavedQuery(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryHandler.RunSavedQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryHandler.RunSavedQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryHandler_fetchQuery(t *testing.T) {
	type fields struct {
		queryStore QueryStore
		osHandler  os.Handler
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qh := Handler{
				queryStore: tt.fields.queryStore,
				osHandler:  tt.fields.osHandler,
			}
			got, err := qh.fetchQuery(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryHandler.fetchQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("QueryHandler.fetchQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryHandler_executeQuery(t *testing.T) {
	type fields struct {
		queryStore QueryStore
		osHandler  os.Handler
	}
	type args struct {
		query string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ResultDTO
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qh := Handler{
				queryStore: tt.fields.queryStore,
				osHandler:  tt.fields.osHandler,
			}
			got, err := qh.executeQuery(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryHandler.executeQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryHandler.executeQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
