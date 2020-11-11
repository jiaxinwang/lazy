package lazy

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.TraceLevel)
}
func Test_separateParams(t *testing.T) {
	type args struct {
		whole map[string][]string
		keys  []string
	}
	tests := []struct {
		name          string
		args          args
		wantSeparated map[string][]string
		wantRemain    map[string][]string
	}{
		{
			"case-1",
			args{whole: map[string][]string{"a": []string{`v1`, `v2`}}, keys: []string{`a`}},
			map[string][]string{"a": []string{`v1`, `v2`}},
			map[string][]string{},
		},
		{
			"case-2",
			args{whole: map[string][]string{"a": []string{`v1`, `v2`}}, keys: []string{}},
			map[string][]string{},
			map[string][]string{"a": []string{`v1`, `v2`}},
		},
		{
			"case-3",
			args{whole: map[string][]string{"a": []string{`v1`, `v2`}, "b": []string{`v3`, `v4`}}, keys: []string{`a`}},
			map[string][]string{"a": []string{`v1`, `v2`}},
			map[string][]string{"b": []string{`v3`, `v4`}},
		},
		{
			"case-4",
			args{whole: map[string][]string{"a": []string{`v1`, `v2`}, "b": []string{`v3`, `v4`}}, keys: []string{`c`}},
			map[string][]string{},
			map[string][]string{"a": []string{`v1`, `v2`}, "b": []string{`v3`, `v4`}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSeparated, gotRemain := separateParams(tt.args.whole, tt.args.keys...)
			if !cmp.Equal(gotSeparated, tt.wantSeparated) {
				t.Errorf("separateParams() gotSeparated = %v, want %v\ndiff=%v", gotSeparated, tt.wantSeparated, cmp.Diff(gotSeparated, tt.wantSeparated))
			}
			if !cmp.Equal(gotRemain, tt.wantRemain) {
				t.Errorf("separateParams() gotRemain = %v, want %v\ndiff=%v", gotRemain, tt.wantRemain, cmp.Diff(gotRemain, tt.wantRemain))
			}
		})
	}
}

func Test_separatePrefixParams(t *testing.T) {
	type args struct {
		whole  map[string][]string
		prefix string
	}
	tests := []struct {
		name          string
		args          args
		wantSeparated map[string][]string
		wantRemain    map[string][]string
	}{
		{
			"case-1",
			args{whole: map[string][]string{"p_a": []string{`v1`, `v2`}}, prefix: `p_`},
			map[string][]string{"p_a": []string{`v1`, `v2`}},
			map[string][]string{},
		},
		{
			"case-2",
			args{whole: map[string][]string{"a": []string{`v1`, `v2`}}, prefix: `p_`},
			map[string][]string{},
			map[string][]string{"a": []string{`v1`, `v2`}},
		},
		{
			"case-3",
			args{whole: map[string][]string{"p_a": []string{`v1`, `v2`}, "b": []string{`v3`, `v4`}}, prefix: `p_`},
			map[string][]string{"p_a": []string{`v1`, `v2`}},
			map[string][]string{"b": []string{`v3`, `v4`}},
		},
		{
			"case-4",
			args{whole: map[string][]string{"p_a": []string{`v1`, `v2`}, "p_b": []string{`v3`, `v4`}}, prefix: `p_`},
			map[string][]string{"p_a": []string{`v1`, `v2`}, "p_b": []string{`v3`, `v4`}},
			map[string][]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSeparated, gotRemain := separatePrefixParams(tt.args.whole, tt.args.prefix)
			if !cmp.Equal(gotSeparated, tt.wantSeparated) {
				t.Errorf("separatePrefixParams() gotSeparated = %v, want %v\ndiff=%v", gotSeparated, tt.wantSeparated, cmp.Diff(gotSeparated, tt.wantSeparated))
			}
			if !cmp.Equal(gotRemain, tt.wantRemain) {
				t.Errorf("separatePrefixParams() gotRemain = %v, want %v\ndiff=%v", gotRemain, tt.wantRemain, cmp.Diff(gotRemain, tt.wantRemain))
			}
		})
	}
}

func Test_separatePage(t *testing.T) {
	type args struct {
		params map[string][]string
	}
	tests := []struct {
		name       string
		args       args
		wantRemain map[string][]string
		wantPage   uint64
		wantLimit  uint64
		wantOffset uint64
	}{
		{`case-1`, args{params: map[string][]string{}}, map[string][]string{}, 0, 0, 0},
		{
			`case-2`,
			args{params: map[string][]string{`offset`: []string{`10`}, `limit`: []string{`2`}, `page`: []string{`3`}}},
			map[string][]string{}, 3, 2, 10,
		},
		{
			`case-3`,
			args{params: map[string][]string{`offset`: []string{`10`, `20`}, `limit`: []string{`2`}, `page`: []string{`3`}}},
			map[string][]string{}, 3, 2, 0,
		},
		{
			`case-4`,
			args{params: map[string][]string{`unused`: []string{`used`}, `offset`: []string{`10`, `20`}, `limit`: []string{`2`}, `page`: []string{`3`}}},
			map[string][]string{`unused`: []string{`used`}}, 3, 2, 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRemain, gotPage, gotLimit, gotOffset := separatePage(tt.args.params)
			if !cmp.Equal(gotRemain, tt.wantRemain) {
				t.Errorf("separatePage() gotRemain = %v, want %v\ndiff=%v", gotRemain, tt.wantRemain, cmp.Diff(gotRemain, tt.wantRemain))
			}
			if gotPage != tt.wantPage {
				t.Errorf("separatePage() gotPage = %v, want %v", gotPage, tt.wantPage)
			}
			if gotLimit != tt.wantLimit {
				t.Errorf("separatePage() gotLimit = %v, want %v", gotLimit, tt.wantLimit)
			}
			if gotOffset != tt.wantOffset {
				t.Errorf("separatePage() gotOffset = %v, want %v", gotOffset, tt.wantOffset)
			}
		})
	}
}

func Test_convertJSONMap(t *testing.T) {
	type args struct {
		src  map[string]interface{}
		maps []JSONPathMap
	}
	tests := []struct {
		name     string
		args     args
		wantDest map[string]interface{}
	}{
		{"case-1", args{map[string]interface{}{"k": "v"}, nil}, map[string]interface{}{"k": "v"}},
		{"case-2", args{map[string]interface{}{"k1": map[string]interface{}{"k2": "v"}}, nil}, map[string]interface{}{"k1": map[string]interface{}{"k2": "v"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDest := SetJSON(tt.args.src, tt.args.maps); !cmp.Equal(gotDest, tt.wantDest) {
				t.Errorf("convertJSONMap() = %v, want %v\ndiff=%v", gotDest, tt.wantDest, cmp.Diff(gotDest, tt.wantDest))
			}
		})
	}
}

func Test_transJSONSingle(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	type args struct {
		src  string
		dest string
		m    JSONPathMap
	}
	tests := []struct {
		name            string
		args            args
		wantConvertSrc  string
		wantConvertDesc string
		wantErr         bool
	}{
		{
			"case-1",
			args{
				`{"name":{"first":"Janet","last":"Prichard"},"age":47}`,
				`{}`,
				JSONPathMap{"age", "age", nil, true},
			},
			`{"name":{"first":"Janet","last":"Prichard"}}`,
			`{"age":47}`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := make(map[string]interface{})
			json.UnmarshalFromString(tt.args.src, &src)
			dest := make(map[string]interface{})
			json.UnmarshalFromString(tt.args.dest, &dest)
			// gotConvertSrc, gotConvertDesc, err := transJSONSingle(src, dest, tt.args.m)
			gotConvertSrc1, gotConvertDesc1, err := SetJSONSingle(src, dest, tt.args.m)
			gotConvertSrc, _ := json.MarshalToString(gotConvertSrc1)
			gotConvertDesc, _ := json.MarshalToString(gotConvertDesc1)
			if (err != nil) != tt.wantErr {
				t.Errorf("transJSONSingle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(gotConvertSrc, tt.wantConvertSrc) {
				t.Errorf("transJSONSingle() gotConvertSrc = %v, want %v\ndiff=%v", gotConvertSrc, tt.wantConvertSrc, cmp.Diff(gotConvertSrc, tt.wantConvertSrc))
			}
			if !cmp.Equal(gotConvertDesc, tt.wantConvertDesc) {
				t.Errorf("transJSONSingle() gotConvertDesc = %v, want %v\ndiff=%v", gotConvertDesc, tt.wantConvertDesc, cmp.Diff(gotConvertDesc, tt.wantConvertDesc))
			}
		})
	}
}
