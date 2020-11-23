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

	config, err := ConfigurationWithContext(c)
	if err != nil {
		return nil, err
	}

	bodyParamsJSONStr, err := json.MarshalToString(bodyParams)
	if err != nil {
		logrus.WithError(err).Error()
		return nil, err
	}

	rels, err := relationships(config.DB, config.Model)
	for _, v := range rels.Relations {
		jsonTag := v.Field.StructField.Tag.Get("json")
		bodyParamsJSONStr, _ = sjson.Delete(bodyParamsJSONStr, jsonTag)

	}

	err = json.UnmarshalFromString(bodyParamsJSONStr, &config.Model)
	if err != nil {
		logrus.WithError(err).Error()
		return nil, err
	}

	newValue := map[string]interface{}{}
	if err := json.UnmarshalFromString(bodyParamsJSONStr, &newValue); err != nil {
		logrus.WithError(err).Error()
		return nil, err
	}
	delete(newValue, "id")

	err = config.DB.Model(config.Model).Updates(newValue).Error
	if err != nil {
		logrus.WithError(err).Error()
		return nil, err
	}

	bodyParamsJSONStr, _ = json.MarshalToString(bodyParams)

	if err == nil {
		for _, v := range rels.Relations {
			jsonTag := v.Field.StructField.Tag.Get("json")
			if _, ok := bodyParams[jsonTag]; ok {
				switch v.Type {
				case schema.Many2Many, schema.HasMany:
					json.UnmarshalFromString(bodyParamsJSONStr, config.Model)
					if fieldValue, err := valueOfField(config.Model, v.Name); err == nil {
						if fieldValue != nil {
							config.DB.Model(config.Model).Association(v.Name).Replace(fieldValue)
						}
					} else {
						logrus.WithError(err).Error()
						return nil, err
					}

					if bodyParamsJSONStr, err = sjson.Delete(bodyParamsJSONStr, v.Field.StructField.Tag.Get("json")); err != nil {
						logrus.WithError(err).Error()
						return nil, err
					}
				case schema.HasOne:
					// do nothing
					if bodyParamsJSONStr, err = sjson.Delete(bodyParamsJSONStr, v.Field.StructField.Tag.Get("json")); err != nil {
						logrus.WithError(err).Error()
						// TODO: error
						return nil, err
					}
				}
			}

		}
	} else {
		logrus.WithError(err).Error()
		return nil, err
	}

	return data, nil
}
