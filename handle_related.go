package lazy

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/schema"
)

// DefaultRelatedPostAction ...
func DefaultRelatedPostAction(c *gin.Context, actionConfig *Action, payload interface{}) (data []map[string]interface{}, err error) {
	urlPaths := strings.Split(c.Request.URL.Path, `/`)
	fieldJSONName := urlPaths[len(urlPaths)-1]

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

	for _, v := range m.Relationships.Relations {
		if strings.EqualFold(v.Field.StructField.Tag.Get("json"), fieldJSONName) {
			switch v.Type {
			case schema.HasMany:
				if fieldValue, err := valueOfField(config.Model, v.Name); err == nil {
					config.DB.Model(config.Model).Association(v.Name).Replace(fieldValue)
				} else {
					logrus.WithError(err).Error()
					return nil, err
				}
			case schema.Many2Many:
				if fieldValue, err := valueOfField(config.Model, v.Name); err == nil {
					s, _ := json.MarshalToString(fieldValue)
					m := []map[string]interface{}{}
					json.UnmarshalFromString(s, &m)

					tagJSON := v.FieldSchema.PrimaryFields[0].StructField.Tag.Get("json")

					ids := []interface{}{}

					for _, vv := range m {
						if value, ok := vv[tagJSON]; ok {
							ids = append(ids, value)
						}
					}
					primaryRets, err := query2Map(config.DB, fmt.Sprintf("%s in ?", tagJSON), ids, v.FieldSchema.Table)
					if err != nil {
						logrus.WithError(err).Error()
						return nil, err
					}

					var inters []interface{}
					for _, vvv := range primaryRets {
						single, _ := NewStruct(v.FieldSchema.ModelType.Name())
						MapStruct(vvv, &single)
						inters = append(inters, single)
					}

					setFieldWithJSONString(config.Model, v.Name, nil)
					config.DB.Model(config.Model).Association(v.Name).Clear()

					if err := setFieldWithJSONString(config.Model, v.Field.StructField.Tag.Get("json"), inters); err != nil {
						logrus.WithError(err).Error()
						return nil, err
					}
					if err := config.DB.Model(config.Model).Association(v.Name).Append(inters); err != nil {
						logrus.WithError(err).Error()
						return nil, err
					}
				} else {
					logrus.WithError(err).Error()
					return nil, err
				}
			case schema.HasOne:
			}
		}
	}
	return
}

// DefaultRelatedPutAction ...
func DefaultRelatedPutAction(c *gin.Context, actionConfig *Action, payload interface{}) (data []map[string]interface{}, err error) {
	return
}

// DefaultRelatedPathcAction ...
func DefaultRelatedPathcAction(c *gin.Context, actionConfig *Action, payload interface{}) (data []map[string]interface{}, err error) {
	return
}

// DefaultRelatedDeleteAction ...
func DefaultRelatedDeleteAction(c *gin.Context, actionConfig *Action, payload interface{}) (data []map[string]interface{}, err error) {
	return
}
