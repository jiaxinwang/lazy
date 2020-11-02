package lazy

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// DefaultDeleteAction ...
func DefaultDeleteAction(c *gin.Context, actionConfig *Action, payload interface{}) (data []map[string]interface{}, err error) {
	config, err := ConfigurationWithContext(c)
	if err != nil {
		return nil, err
	}
	contextParams, err := params(c)
	if err != nil {
		return nil, ErrParamMissing
	}

	ids := valueSliceWithParamKey(contextParams, "id")

	m, err := schema.Parse(config.Model, schemaStore, schema.NamingStrategy{})

	for _, vID := range ids {
		priDBName := m.PrimaryFieldDBNames[0]
		mapResult := map[string]interface{}{}
		// TODO: batch
		err := config.DB.Model(config.Model).Where(fmt.Sprintf("%s = ?", priDBName), vID).Find(&mapResult).Error
		if err != nil {
			logrus.WithError(err).Warn()
			continue
		}
		cloned := clone(config.Model)
		MapStruct(mapResult, cloned)

		for _, v := range m.Relationships.Relations {
			err = config.DB.Model(cloned).Association(v.Name).Clear()
			if err != nil && err != gorm.ErrRecordNotFound {
				logrus.WithError(err).Error()
				return nil, err
			}
		}
	}

	err = config.DB.Model(config.Model).Where(fmt.Sprintf("%s in ?", m.PrimaryFieldDBNames[0]), ids).Delete(config.Model).Error

	return nil, err

}
