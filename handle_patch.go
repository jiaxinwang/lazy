package lazy

import (
	"reflect"
	"strconv"
	"strings"
	"time"

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
		if !strings.EqualFold(jsonTag, "") {
			var deleteErr error
			bodyParamsJSONStr, deleteErr = sjson.Delete(bodyParamsJSONStr, jsonTag)
			if deleteErr != nil {
				logrus.WithError(err).Error()
				return nil, err
			}
		}
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

	// vvv := newValue["birth"]
	// logrus.Infof("%#v", vvv)
	// logrus.Panicf("%+v", newValue)

	m, err := schema.Parse(config.Model, schemaStore, schema.NamingStrategy{})
	if err != nil {
		logrus.WithError(err).Errorf("can't parse %+v", config.Model)
	} else {
		for k, v := range newValue {
			if f, ok := m.FieldsByDBName[k]; ok {
				switch f.StructField.Type.Kind() {
				case reflect.TypeOf(time.Time{}).Kind(), reflect.TypeOf(&time.Time{}).Kind():
					str := v.(string)
					if t, err := time.Parse(time.RFC3339, str); err == nil {
						newValue[k] = t
					} else {
						logrus.WithError(err).Error()
					}
					// logrus.Printf("%#v", newValue[k])
				}
			}
		}
	}

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
