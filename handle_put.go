package lazy

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// DefaultPutAction ...
func DefaultPutAction(c *gin.Context, actionConfig *Action, payload interface{}) (data []map[string]interface{}, err error) {
	_, queryParams, bodyParams := ContentParams(c)
	ids, ok := queryParams[`id`]
	mids := ids.([]string)
	if !ok || len(mids) != 1 {
		return nil, ErrParamMissing
	}

	id, err := strconv.Atoi(mids[0])
	if err != nil {
		return nil, err
	}
	bodyParams[`id`] = id

	config, err := ConfigurationWithContext(c)
	if err != nil {
		return nil, err
	}
	s, err := json.MarshalToString(bodyParams)
	if err != nil {
		logrus.WithError(err).Error()
		// TODO: error
		return nil, err
	}
	err = json.UnmarshalFromString(s, &config.Model)
	if err != nil {
		logrus.WithError(err).Error()
		// TODO: error
		return nil, err
	}

	m, err := schema.Parse(config.Model, schemaStore, schema.NamingStrategy{})

	priDBName := m.PrimaryFieldDBNames[0]
	mapResult := map[string]interface{}{}
	err = config.DB.Model(config.Model).Where(fmt.Sprintf("%s = ?", priDBName), id).Find(&mapResult).Error
	if err != nil {
		logrus.WithError(err).Error()
		return nil, err
	}
	cloned := clone(config.Model)
	MapStruct(mapResult, cloned)

	for _, v := range m.Relationships.HasMany {
		err = config.DB.Model(cloned).Association(v.Name).Clear()
		if err != nil && err != gorm.ErrRecordNotFound {
			logrus.WithError(err).Error()
			return nil, err
		}
	}
	for _, v := range m.Relationships.Many2Many {
		err = config.DB.Model(cloned).Association(v.Name).Clear()
		if err != nil && err != gorm.ErrRecordNotFound {
			logrus.WithError(err).Error()
			return nil, err
		}
	}

	updateModel(config.DB, config.Model)

	return
}
