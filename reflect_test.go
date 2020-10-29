package lazy

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

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
