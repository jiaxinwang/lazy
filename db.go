package lazy

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/schema"
)

type disassembled struct {
	Schema     string
	StructName string
	Values     []interface{}
}

func disassemble(db *gorm.DB, model interface{}) (many2many []disassembled, err error) {
	var b []byte
	if b, err = json.Marshal(model); err != nil {
		return
	}

	sfs := db.NewScope(model).GetStructFields()
	for _, v := range sfs {
		if r := v.Relationship; r != nil {
			switch r.Kind {
			case string(schema.BelongsTo):
			case string(schema.HasOne):
			case string(schema.HasMany):
			case string(schema.Many2Many):
				m2m := disassembled{Schema: string(schema.Many2Many), StructName: v.Name}

				sub := json.Get(b, v.Struct.Tag.Get("json"))
				associated := make([]map[string]interface{}, 0)
				json.Unmarshal([]byte(sub.ToString()), &associated)
				m2m.Values = make([]interface{}, len(associated))
				for k, vv := range associated {
					if len(r.ForeignFieldNames) > 0 {
						s, ok := newStruct(v.Struct.Type.Elem().Name())
						if !ok {
							continue
						}
						if err = MapStruct(vv, &s); err != nil {
							return
						}
						m2m.Values[k] = s
					}
				}
				many2many = append(many2many, m2m)
			}
		}
	}
	return
}

func createModel(db *gorm.DB, model interface{}) (err error) {
	var b []byte
	if b, err = json.Marshal(model); err != nil {
		return
	}
	if err = db.Create(model).Error; err != nil {
		return
	}
	sfs := db.NewScope(model).GetStructFields()
	for _, v := range sfs {
		if v.Relationship != nil {
			r := v.Relationship
			switch r.Kind {
			case string(schema.BelongsTo):
				logrus.WithField("kind", r.Kind).Tracef("do nothing")
			case string(schema.HasOne):
				logrus.WithField("kind", r.Kind).Tracef("do nothing")
			case string(schema.HasMany):
				logrus.WithField("kind", r.Kind).Tracef("do nothing")
			case string(schema.Many2Many):
				r := v.Relationship
				sub := json.Get(b, v.Struct.Tag.Get("json"))
				associated := make([]map[string]interface{}, 0)
				json.Unmarshal([]byte(sub.ToString()), &associated)
				for _, vv := range associated {
					if len(r.ForeignFieldNames) > 0 {
						s, ok := newStruct(v.Struct.Type.Elem().Name())
						if !ok {
							continue
						}
						if err = MapStruct(vv, &s); err != nil {
							return err
						}
						logrus.Print(111222333)
						if err := db.Model(model).Association(v.Name).Append(s).Error; err != nil {
							logrus.WithError(err).Error()
							return err
						}
					}
				}

			}
		}
	}

	return
}
