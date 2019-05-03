package query

import (
	"reflect"
	"testing"
)

func TestNewLocalQueryStore(t *testing.T) {
	tests := []struct {
		name    string
		queries []map[string]string
		want    *localQueryStore
	}{
		{
			name:    "No standard queries",
			queries: nil,
			want:    &localQueryStore{queries: map[string]string{}},
		},
		{
			name:    "One map of standard queries",
			queries: []map[string]string{
				{
					"foo": "bar",
				},
			},
			want:    &localQueryStore{queries: map[string]string{"foo": "bar",}},
		},
		{
			name:    "Multiple maps of standard queries",
			queries: []map[string]string{
				{
					"foo1": "bar1",
				},
				{
					"foo2": "bar2",
					"foo3": "bar3",
				},
				{
					"foo4": "bar4",
					"foo5": "bar5",
					"foo6": "bar6",
				},
			},
			want:    &localQueryStore{queries: map[string]string{
					"foo1": "bar1",
					"foo2": "bar2",
					"foo3": "bar3",
					"foo4": "bar4",
					"foo5": "bar5",
					"foo6": "bar6",
				},
			},
		},
	}
	for _, tt := range tests {
		// override the standard queries for the test
		queries = tt.queries

		t.Run(tt.name, func(t *testing.T) {
			if got := NewLocalQueryStore(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLocalQueryStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_localQueryStore_GetQuery(t *testing.T) {
	type fields struct {
		queries map[string]string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  bool
	}{
		{
			name: "Nil queries, empty name",
			fields: fields{
				queries: nil,
			},
			args:  args{name: ""},
			want:  "",
			want1: false,
		},
		{
			name: "Nil queries, non-empty name",
			fields: fields{
				queries: nil,
			},
			args:  args{name: "users"},
			want:  "",
			want1: false,
		},
		{
			name: "Initialised queries, empty name",
			fields: fields{
				queries: map[string]string{
					"query": "test",
				},
			},
			args:  args{name: ""},
			want:  "",
			want1: false,
		},
		{
			name: "Nil queries, query not in store",
			fields: fields{
				queries: map[string]string{
					"query": "test",
				},
			},
			args:  args{name: "other_query"},
			want:  "",
			want1: false,
		},
		{
			name: "Initialised queries, query in store",
			fields: fields{
				queries: map[string]string{
					"query": "test",
				},
			},
			args:  args{name: "query"},
			want:  "test",
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qs := &localQueryStore{
				queries: tt.fields.queries,
			}
			got, got1 := qs.GetQuery(tt.args.name)
			if got != tt.want {
				t.Errorf("localQueryStore.GetQuery() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("localQueryStore.GetQuery() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

