package query

import (
	"testing"
)

/*
func TestGetQuery(t *testing.T) {
	type args struct {
		option string
	}
	tests := []struct {
		query  string
		args  args
		want  string
		want1 bool
	}{
		0: {
			query:  "has kernel_info",
			args:  args{option: "kernel_info"},
			want:  "SELECT * FROM kernel_info;",
			want1: true,
		},
		1: {
			query:  "typo in kernel_info",
			args:  args{option: "kerel_info"},
			want:  "",
			want1: false,
		},
		2: {
			query:  "unknown query",
			args:  args{option: "foobar"},
			want:  "",
			want1: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
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
			fields:fields{
				queries: nil,
			},
			args: args{name: ""},
			want: "",
			want1: false,
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

