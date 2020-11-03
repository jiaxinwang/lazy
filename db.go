package lazy

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func createModel(db *gorm.DB, model interface{}) (err error) {
	return db.Create(model).Error
}

func updateModel(db *gorm.DB, model interface{}) (err error) {
	return db.Save(model).Error
}

func relationships(db *gorm.DB, model interface{}) (relationships schema.Relationships, err error) {
	m, err := schema.Parse(model, schemaStore, schema.NamingStrategy{})
	if err != nil {
		return schema.Relationships{}, err
	}

	return m.Relationships, nil
}

func query2Map(db *gorm.DB, condition string, where []interface{}, table string) (result []map[string]interface{}, err error) {
	result = []map[string]interface{}{}
	if !strings.EqualFold(condition, "") {
		return result, db.Table(table).Where(condition, where).Find(&result).Error
	}
	return result, db.Table(table).Find(&result).Error
}
