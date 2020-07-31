package lazy

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/schema"
)

type disassembled struct {
	DBName            string
	Schema            string
	StructName        string
	TypeName          string
	ForeignDBNames    []string
	ForeignFieldNames []string
	Values            []interface{}
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
				m2m := disassembled{DBName: v.DBName, Schema: string(schema.Many2Many), TypeName: v.Struct.Type.Elem().Name(), StructName: v.Struct.Name}
				m2m.ForeignDBNames = r.ForeignDBNames
				m2m.ForeignFieldNames = r.ForeignFieldNames
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
	m2m := make([]disassembled, 0)
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
			}
		}
	}

	if m2m, err = disassemble(db, model); err == nil {
		for _, v := range m2m {
			if err = setField(model, v.StructName, nil); err != nil {
				return err
			}
		}
	} else {
		return err
	}

	if err = db.Create(model).Error; err != nil {
		return
	}

	for _, v := range m2m {
		if len(v.Values) == 0 {
			continue
		}
		if f, ok := newStruct(v.TypeName); ok {
			m := make([]interface{}, 0)
			eqs := make(map[string][]interface{})
			for _, vvv := range v.Values {
				key := v.ForeignFieldNames[0]
				m = append(m, valueOfJSONKey(vvv, key).ToString())
			}
			eqs[v.ForeignFieldNames[0]] = m

			sel := sq.Select("*").From(v.DBName)
			sel = SelectBuilder(sel, eqs, nil, nil, nil, nil)
			data, _ := ExecSelect(db, sel)

			for _, vv := range data {
				if err = MapStruct(vv, &f); err != nil {
					return err
				}
				db.Model(model).Association(v.StructName).Append(f)
			}
		} else {
			return
		}
	}
	return
}
