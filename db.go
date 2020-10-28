package lazy

import (
	"fmt"
	"strings"
	"sync"

	sq "github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func createModel(db *gorm.DB, model interface{}) (err error) {

	return db.Create(model).Error
}

func updateModel(db *gorm.DB, model interface{}) (err error) {
	return nil
}

func relationships(db *gorm.DB, model interface{}) (relationships schema.Relationships, err error) {
	m, err := schema.Parse(model, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		return schema.Relationships{}, err
	}

	return m.Relationships, nil
}

func associateModel(db *gorm.DB, model interface{}) (err error) {
	m, err := schema.Parse(model, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		return err
	}
	// for _, v := range m.Relationships.BelongsTo {
	// 	logrus.Infof("%#v", v)
	// }
	for _, v := range m.Relationships.Many2Many {
		logrus.Infof("1 %#v", v)
		logrus.Infof("2 %#v", v.References[0])
		logrus.Infof("3 %#v", v.References[0].ForeignKey)
		logrus.Infof("4 %#v", v.References[0].PrimaryValue)
		logrus.Infof("5 %#v", v.References[0].PrimaryKey)
		logrus.Infof("6 %#v", v.Field.Name)

		// logrus.Print(schema.NamingStrategy{}.ColumnName("Foods", "ID"))

		joined := strings.Join([]string{m.Name, v.Name}, "")
		tableName := schema.NamingStrategy{}.TableName(joined)
		logrus.Print(tableName)

		// if part := strings.Split(v.Field.StructField.Type.Elem().String(), `.`); len(part) > 0 {
		// joined := strings.Join([]string{v.Name}, ".")
		// tableName := schema.NamingStrategy{}.TableName(part[len(part)-1])
		// tableName := schema.NamingStrategy{}.TableName(joined)
		// logrus.Print(tableName)
		primaryValue, err := valueOfField(model, v.References[0].PrimaryKey.Name)
		if err != nil {
			return fmt.Errorf("associate model %v: %w", model, err)
		}
		var results []map[string]interface{}
		db.Table(tableName).Where(fmt.Sprintf("%s = ?", v.References[0].ForeignKey.DBName), primaryValue).Find(&results)
		logrus.Print(results)
		// set := make([]interface{}, 0)
		// for _, vv := range results {
		// 	f, ok := NewStruct(part[len(part)-1])
		// 	if !ok {
		// 		return fmt.Errorf("associate model new struct %v %v", model, part[len(part)-1])
		// 	}
		// 	str, err := json.MarshalToString(vv)
		// 	if err != nil {
		// 		return fmt.Errorf("associate model marshal to string %v: %w", model, err)
		// 	}

		// 	err = json.UnmarshalFromString(str, &f)
		// 	if err != nil {
		// 		return fmt.Errorf("associate model unmarshal from string %v: %w", model, err)
		// 	}

		// 	set = append(set, f)
		// }
		// err = setFieldWithJSONString(model, v.Field.StructField.Tag.Get("json"), set)
		// if err != nil {
		// 	return fmt.Errorf("associate model set field by json %v (%s): %w", model, v.Field.StructField.Tag.Get("json"), err)
		// }
		// }

	}
	for _, v := range m.Relationships.HasMany {
		if part := strings.Split(v.Field.StructField.Type.Elem().String(), `.`); len(part) > 0 {
			tableName := schema.NamingStrategy{}.TableName(part[len(part)-1])

			primaryValue, err := valueOfField(model, v.References[0].PrimaryKey.Name)
			if err != nil {
				return fmt.Errorf("associate model %v: %w", model, err)
			}
			var results []map[string]interface{}
			db.Table(tableName).Where(fmt.Sprintf("%s = ?", v.References[0].ForeignKey.DBName), primaryValue).Find(&results)
			set := make([]interface{}, 0)
			for _, vv := range results {
				f, ok := NewStruct(part[len(part)-1])
				if !ok {
					return fmt.Errorf("associate model new struct %v %v", model, part[len(part)-1])
				}
				str, err := json.MarshalToString(vv)
				if err != nil {
					return fmt.Errorf("associate model marshal to string %v: %w", model, err)
				}

				err = json.UnmarshalFromString(str, &f)
				if err != nil {
					return fmt.Errorf("associate model unmarshal from string %v: %w", model, err)
				}

				set = append(set, f)
			}
			err = setFieldWithJSONString(model, v.Field.StructField.Tag.Get("json"), set)
			if err != nil {
				return fmt.Errorf("associate model set field by json %v (%s): %w", model, v.Field.StructField.Tag.Get("json"), err)
			}

		}
	}

	// for k, v := range m.Relationships.HasOne {
	// 	logrus.WithField("k", k).WithField("v", fmt.Sprintf("%#v", v)).Info()
	// }

	// sfs := db.NewScope(model).GetStructFields()
	// for _, v := range sfs {
	// 	if v.Relationship != nil {
	// 		r := v.Relationship
	// 		switch r.Kind {
	// 		case string(schema.BelongsTo):
	// 			// logrus.Panic()
	// 		case string(schema.HasOne):
	// 			av, err := valueOfField(model, r.AssociationForeignFieldNames[0])
	// 			if err != nil {
	// 				return err
	// 			}
	// 			if f, ok := newStruct(v.Struct.Type.Name()); ok {
	// 				sel := sq.Select("*").From(db.NewScope(f).TableName())
	// 				sel = SelectBuilder(sel, map[string][]interface{}{r.ForeignDBNames[0]: {av}}, nil, nil, nil, nil)
	// 				data, _ := ExecSelect(db, sel)
	// 				if err = MapStruct(data[0], &f); err != nil {
	// 					return err
	// 				}
	// 				setField(model, v.Name, f)
	// 			} else {
	// 				return ErrHasAssociations
	// 			}
	// 		case string(schema.HasMany):
	// 			av, err := valueOfField(model, r.AssociationForeignFieldNames[0])
	// 			if err != nil {
	// 				return err
	// 			}
	// 			set := make([]interface{}, 0)
	// 			if f, ok := newStruct(v.Struct.Type.Elem().Name()); ok {
	// 				sel := sq.Select("*").From(db.NewScope(f).TableName())
	// 				sel = SelectBuilder(sel, map[string][]interface{}{r.ForeignDBNames[0]: {av}}, nil, nil, nil, nil)
	// 				data, _ := ExecSelect(db, sel)
	// 				for _, datav := range data {
	// 					if err = MapStruct(datav, &f); err != nil {
	// 						return err
	// 					}
	// 					set = append(set, f)

	// 				}
	// 			} else {
	// 				return ErrHasAssociations
	// 			}
	// 			setJSONField(model, v.Tag.Get("json"), set)
	// 		case string(schema.Many2Many):
	// 			aaa := valueOfJSONKey(model, r.AssociationForeignFieldNames[0]).ToString()
	// 			cond := fmt.Sprintf(`"%s"."%s" = %s`, r.JoinTableHandler.Table(db), r.ForeignDBNames[0], aaa)
	// 			join := fmt.Sprintf(
	// 				`INNER JOIN %s ON "%s"."%s" = "%s"."%s"`,
	// 				r.JoinTableHandler.Table(db),
	// 				r.JoinTableHandler.Table(db),
	// 				r.AssociationForeignDBNames[0],
	// 				v.DBName,
	// 				r.AssociationForeignFieldNames[0],
	// 			)

	// 			sel := sq.Select("*").From(v.DBName).JoinClause(join).Where(cond)

	// 			data, _ := ExecSelect(db, sel)
	// 			if err = setJSONField(model, v.Tag.Get("json"), data); err != nil {
	// 				logrus.WithError(err).Error()
	// 				return err
	// 			}
	// 		}
	// 	}
	// }
	return nil
}

func queryAssociated(db *gorm.DB, foreignDBName, foreignFieldName string, foreignFieldValue interface{}) (ret []map[string]interface{}) {
	eqs := make(map[string][]interface{})
	eqs[foreignFieldName] = []interface{}{foreignFieldValue}
	sel := sq.Select("*").From(foreignDBName)
	sel = SelectBuilder(sel, eqs, nil, nil, nil, nil)
	ret, _ = ExecSelect(db, sel)
	return
}
