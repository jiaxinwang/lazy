package lazy

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/sjson"
	"gorm.io/gorm/schema"
)

// DefaultPatchAction ...
func DefaultPatchAction(c *gin.Context, actionConfig *Action, payload interface{}) (data []map[string]interface{}, err error) {
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

	// ids :=

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

	rels, err := relationships(config.DB, config.Model)
	if err == nil {
		for _, v := range rels.Relations {
			switch v.Type {
			case schema.HasMany, schema.HasOne, schema.Many2Many:
				logrus.Print(v.Field.StructField.Tag.Get("json"))
				logrus.Print(s)
				if s, err = sjson.Delete(s, v.Field.StructField.Tag.Get("json")); err != nil {
					logrus.WithError(err).Error()
					// TODO: error
					return nil, err
				}
				logrus.Print(s)
			}
		}
	} else {
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

	newValue := map[string]interface{}{}
	if err := json.UnmarshalFromString(s, &newValue); err != nil {
		logrus.WithError(err).Error()
		// TODO: error
		return nil, err
	}

	err = config.DB.Model(config.Model).Updates(newValue).Error
	if err != nil {
		logrus.WithError(err).Error()
		// TODO: error
		return nil, err
	}

	// TODO: return
	return data, nil
}
