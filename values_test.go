package lazy

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_valueOfMap(t *testing.T) {
	map1 := map[string][]string{
		"limit":  {"10"},
		"offset": {"500"},
		"page":   {"2"},
		"body":   {"1", "2", "3"},
	}

	type args struct {
		params map[string][]string
		key    string
	}
	tests := []struct {
		name      string
		args      args
		wantValue []string
		wantOk    bool
	}{
		{"limit", args{map1, "limit"}, []string{"10"}, true},
		{"offset", args{map1, "offset"}, []string{"500"}, true},
		{"page", args{map1, "page"}, []string{"2"}, true},
		{"body", args{map1, "body"}, []string{"1", "2", "3"}, true},
		{"none", args{map1, "none"}, []string{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotOk := valueOfMap(tt.args.params, tt.args.key)
			if !cmp.Equal(gotValue, tt.wantValue) {
				t.Errorf("valueOfMap() gotValue = %v, want %v\ndiff=%v", gotValue, tt.wantValue, cmp.Diff(gotValue, tt.wantValue))
			}
			if gotOk != tt.wantOk {
				t.Errorf("valueOfMap() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
