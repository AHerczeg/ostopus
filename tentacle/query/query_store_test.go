package query

import (
	"reflect"
	"testing"
)

/*
func TestGetQuery(t *testing.T) {
	type args struct {
		option string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		0: {
			name:  "has kernel_info",
			args:  args{option: "kernel_info"},
			want:  "SELECT * FROM kernel_info;",
			want1: true,
		},
		1: {
			name:  "typo in kernel_info",
			args:  args{option: "kerel_info"},
			want:  "",
			want1: false,
		},
		2: {
			name:  "unknown query",
			args:  args{option: "foobar"},
			want:  "",
			want1: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetQuery(tt.args.option)
			if got != tt.want {
				t.Errorf("GetQuery() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetQuery() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
*/

func TestNewQueryStore(t *testing.T) {
	tests := []struct {
		name string
		want localQueryStore
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQueryStore(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQueryStore() = %v, want %v", got, tt.want)
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qs := &localQueryStore{
				queries: tt.fields.queries,
			}
			qs.AddQueries(tt.args.queries)
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
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qs := &localQueryStore{
				queries: tt.fields.queries,
			}
			qs.AddQuery(tt.args.name, tt.args.query)
		})
	}
}

func TestFoo(t *testing.T) {
	type args struct {
		i int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Foo(tt.args.i); got != tt.want {
				t.Errorf("Foo() = %v, want %v", got, tt.want)
			}
		})
	}
}
