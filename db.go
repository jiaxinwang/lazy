package lazy

import (
	sq "github.com/Masterminds/squirrel"
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
	m, err := schema.Parse(model, schemaStore, schema.NamingStrategy{})
	if err != nil {
		return schema.Relationships{}, err
	}

	return m.Relationships, nil
}

func queryAssociated(db *gorm.DB, foreignDBName, foreignFieldName string, foreignFieldValue interface{}) (ret []map[string]interface{}) {
	eqs := make(map[string][]interface{})
	eqs[foreignFieldName] = []interface{}{foreignFieldValue}
	sel := sq.Select("*").From(foreignDBName)
	sel = SelectBuilder(sel, eqs, nil, nil, nil, nil)
	ret, _ = ExecSelect(db, sel)
	return
}
