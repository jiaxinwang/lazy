package lazy

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_valueOfTag(t *testing.T) {
	type args struct {
		inter   interface{}
		tagName string
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{"profile", args{inter: &Profile{}, tagName: "id"}, uint(0)},
		{"profile", args{inter: &Profile{Age: 2}, tagName: "age"}, uint(2)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := valueOfTag(tt.args.inter, tt.args.tagName); !cmp.Equal(got, tt.want) {
				t.Errorf("valueOfTag() = %v, want %v\ndiff=%v", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}
