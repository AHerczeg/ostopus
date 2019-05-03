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
			name: "One map of standard queries",
			queries: []map[string]string{
				{
					"foo": "bar",
				},
			},
			want: &localQueryStore{queries: map[string]string{"foo": "bar"}},
		},
		{
			name: "Multiple maps of standard queries",
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
			want: &localQueryStore{queries: map[string]string{
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

func Test_localQueryStore_HasQuery(t *testing.T) {
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
		want   bool
	}{
		{
			name: "Nil query",
			fields: fields{
				queries: nil,
			},
			args: args{
				name: "",
			},
			want: false,
		},
		{
			name: "Initialised queries, query not in store",
			fields: fields{
				queries: map[string]string{
					"foo": "bar",
				},
			},
			args: args{
				name: "other_query",
			},
			want: false,
		},
		{
			name: "Initialised queries, query in store",
			fields: fields{
				queries: map[string]string{
					"foo": "bar",
				},
			},
			args: args{
				name: "foo",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qs := &localQueryStore{
				queries: tt.fields.queries,
			}
			if got := qs.HasQuery(tt.args.name); got != tt.want {
				t.Errorf("localQueryStore.HasQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_localQueryStore_AddQueries(t *testing.T) {
	type fields struct {
		queries map[string]string
	}
	type args struct {
		queries map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]string
	}{
		{
			name: "Adding nil to nil",
			fields: fields{
				queries: nil,
			},
			args: args{
				queries: nil,
			},
			want: nil,
		},
		{
			name: "Adding nil to initialised queries",
			fields: fields{
				queries: map[string]string{
					"foo": "bar",
				},
			},
			args: args{
				queries: nil,
			},
			want: map[string]string{"foo": "bar"},
		},
		{
			name: "Adding queries to nil",
			fields: fields{
				queries: nil,
			},
			args: args{
				queries: map[string]string{
					"foo": "bar",
				},
			},
			want: nil,
		},
		{
			name: "Adding one query",
			fields: fields{
				queries: map[string]string{
					"foo": "bar",
				},
			},
			args: args{
				queries: map[string]string{
					"new": "query",
				},
			},
			want: map[string]string{"foo": "bar", "new": "query"},
		},
		{
			name: "Adding multiple queries",
			fields: fields{
				queries: map[string]string{
					"foo": "bar",
				},
			},
			args: args{
				queries: map[string]string{
					"new1": "query1",
					"new2": "query2",
					"new3": "query3",
				},
			},
			want: map[string]string{
				"foo":  "bar",
				"new1": "query1",
				"new2": "query2",
				"new3": "query3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qs := &localQueryStore{
				queries: tt.fields.queries,
			}
			qs.AddQueries(tt.args.queries)
			if !reflect.DeepEqual(qs.queries, tt.want) {
				t.Errorf("localQueryStore.AddQueris() = %v, want %v", qs.queries, tt.want)
			}
		})
	}
}

func Test_localQueryStore_AddQuery(t *testing.T) {
	type fields struct {
		queries map[string]string
	}
	type args struct {
		name  string
		query string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]string
	}{
		{
			name: "Adding query to nil",
			fields: fields{
				queries: nil,
			},
			args: args{
				name: "foo",
				query: "bar",
			},
			want: nil,
		},
		{
			name: "Adding one query",
			fields: fields{
				queries: map[string]string{
					"foo": "bar",
				},
			},
			args: args{
				name: "new",
				query: "query",
			},
			want: map[string]string{"foo": "bar", "new": "query"},
		},
		{
			name: "Adding one query already in store",
			fields: fields{
				queries: map[string]string{
					"foo": "bar",
				},
			},
			args: args{
				name: "foo",
				query: "bar",
			},
			want: map[string]string{"foo": "bar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qs := &localQueryStore{
				queries: tt.fields.queries,
			}
			qs.AddQuery(tt.args.name, tt.args.query)
			if !reflect.DeepEqual(qs.queries, tt.want) {
				t.Errorf("localQueryStore.AddQuery() = %v, want %v", qs.queries, tt.want)
			}
		})
	}
}


