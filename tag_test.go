package lazy

import "testing"

// func Test_disassembleTag(t *testing.T) {
// 	type args struct {
// 		tag string
// 	}
// 	tests := []struct {
// 		name                string
// 		args                args
// 		wantName            string
// 		wantForeignkeyTable string
// 		wantForeignkey      string
// 		wantErr             bool
// 	}{
// 		{"name", args{tag: "tag_name"}, `tag_name`, ``, ``, false},
// 		{"foreign", args{tag: "tag_name;foreign:table.id"}, `tag_name`, `table`, `id`, false},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			gotName, gotForeignkeyTable, gotForeignkey, err := disassembleTag(tt.args.tag)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("disassembleTag() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if gotName != tt.wantName {
// 				t.Errorf("disassembleTag() gotName = %v, want %v", gotName, tt.wantName)
// 			}
// 			if gotForeignkeyTable != tt.wantForeignkeyTable {
// 				t.Errorf("disassembleTag() gotForeignkeyTable = %v, want %v", gotForeignkeyTable, tt.wantForeignkeyTable)
// 			}
// 			if gotForeignkey != tt.wantForeignkey {
// 				t.Errorf("disassembleTag() gotForeignkey = %v, want %v", gotForeignkey, tt.wantForeignkey)
// 			}
// 		})
// 	}
// }

func Test_disassembleTag(t *testing.T) {
	type args struct {
		tag string
	}
	tests := []struct {
		name                string
		args                args
		wantName            string
		wantId              string
		wantForeignkeyTable string
		wantForeignkey      string
		wantErr             bool
	}{
		{"name", args{tag: "tag_name"}, `tag_name`, ``, ``, ``, false},
		{"foreign", args{tag: "tag_name;foreign:id->table.pid"}, `tag_name`, `id`, `table`, `pid`, false},
		{"dog", args{tag: "profile;foreign:id->profiles.dog_id"}, `profile`, `id`, `profiles`, `dog_id`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotId, gotForeignkeyTable, gotForeignkey, err := disassembleTag(tt.args.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("disassembleTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotName != tt.wantName {
				t.Errorf("disassembleTag() gotName = %v, want %v", gotName, tt.wantName)
			}
			if gotId != tt.wantId {
				t.Errorf("disassembleTag() gotId = %v, want %v", gotId, tt.wantId)
			}
			if gotForeignkeyTable != tt.wantForeignkeyTable {
				t.Errorf("disassembleTag() gotForeignkeyTable = %v, want %v", gotForeignkeyTable, tt.wantForeignkeyTable)
			}
			if gotForeignkey != tt.wantForeignkey {
				t.Errorf("disassembleTag() gotForeignkey = %v, want %v", gotForeignkey, tt.wantForeignkey)
			}
		})
	}
}
