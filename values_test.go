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

func Test_splitParams1(t *testing.T) {
	type args struct {
		model  interface{}
		params map[string][]string
	}
	// zeroQueryParam := QueryParam{}
	tests := []struct {
		name           string
		args           args
		wantQueryParam QueryParam
		wantErr        bool
	}{
		// TODO: Add test cases.
		{"1", args{&Dog{},
			map[string][]string{
				"none":           {"none"},
				"name":           {"name"},
				"name_like":      {"%%"},
				"created_at_lt":  {"2020-01-01 09:00:00"},
				"created_at_gt":  {"2020-01-01 10:00:00"},
				"created_at_lte": {"2020-01-01 11:00:00"},
				"created_at_gte": {"2020-01-01 12:00:00"},
				"foods":          {"1"},
				"toys":           {"2"},
			}},
			QueryParam{
				Eq:        map[string][]interface{}{"name": {[]string{"name"}}},
				Lt:        map[string][]interface{}{"created_at": {[]string{"2020-01-01 09:00:00"}}},
				Lte:       map[string][]interface{}{"created_at": {[]string{"2020-01-01 11:00:00"}}},
				Gt:        map[string][]interface{}{"created_at": {[]string{"2020-01-01 10:00:00"}}},
				Gte:       map[string][]interface{}{"created_at": {[]string{"2020-01-01 12:00:00"}}},
				Like:      map[string][]interface{}{"name": {[]string{"%%"}}},
				HasMany:   map[string][]interface{}{"toys": {[]string{"2"}}},
				Many2Many: map[string][]interface{}{"foods": {[]string{"1"}}},
			}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQueryParam, err := splitParams1(tt.args.model, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("splitParams1() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(gotQueryParam, tt.wantQueryParam) {
				t.Errorf("splitParams1() = %v, want %v\ndiff=%v", gotQueryParam, tt.wantQueryParam, cmp.Diff(gotQueryParam, tt.wantQueryParam))
			}
		})
	}
}
