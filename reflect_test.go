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

func Test_setField(t *testing.T) {
	type args struct {
		src  interface{}
		name string
		v    interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantRet interface{}
		wantErr bool
	}{
		{"case-0", args{src: &Dog{Foods: []Food{{Brand: "1"}}}, name: "Foods", v: []Food{}}, &Dog{Foods: []Food{}}, false},
		{"case-1", args{src: &Dog{Foods: []Food{{Brand: "1"}}}, name: "Foods", v: nil}, &Dog{}, false},
		{"case-2", args{src: &Dog{Name: "good dog"}, name: "Name", v: "bad dog"}, &Dog{Name: "bad dog"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := setField(tt.args.src, tt.args.name, tt.args.v)
			gotRet := tt.args.src
			if (err != nil) != tt.wantErr {
				t.Errorf("setField() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(gotRet, tt.wantRet) {
				t.Errorf("setField() = %v, want %v\ndiff=%v", gotRet, tt.wantRet, cmp.Diff(gotRet, tt.wantRet))
			}
		})
	}
}

func TestParseDescribe(t *testing.T) {
	type args struct {
		inter interface{}
	}

	describe := &Describe{
		Name: "Food",
		Fields: []*Field{
			{Name: "ID"},
			{Name: "CreatedAt"},
			{Name: "UpdatedAt"},
			{Name: "Brand"},
		},
	}

	tests := []struct {
		name    string
		args    args
		want    *Describe
		wantErr bool
	}{
		{"case-1", args{&Food{}}, describe, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDescribe(tt.args.inter)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDescribe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for k, v := range got.Fields {
				if !cmp.Equal(v.Name, tt.want.Fields[k].Name) {
					t.Errorf("ParseDescribe() = %v, want %v\ndiff=%v", v.Name, tt.want.Fields[k].Name, cmp.Diff(v.Name, tt.want.Fields[k].Name))
				}

			}
			// if !cmp.Equal(got, tt.want) {
			// 	t.Errorf("ParseDescribe() = %v, want %v\ndiff=%v", got, tt.want, cmp.Diff(got, tt.want))
			// }
		})
	}
}
