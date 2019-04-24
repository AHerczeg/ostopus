package query

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
