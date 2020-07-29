package lazy

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/schema"
)

func createModel(db *gorm.DB, model interface{}) (err error) {
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
				logrus.WithField("v.Struct.Name", v.Struct.Name).Debug()
				// case string(schema.HasOne), string(schema.Many2Many):
				// 	return nil, ErrUnknown
				// case string(schema.HasMany):
				// 	count := 0
				// 	config.DB.Table(v.DBName).Where(fmt.Sprintf("%s = ?", r.ForeignDBNames[0]), id).Count(&count)
				// 	if count > 0 {
				// 		return nil, ErrHasAssociations
				// 	}
			}
		}
	}

	return
}
