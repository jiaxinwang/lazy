package lazy

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/schema"
)

// DefaultDeleteAction ...
func DefaultDeleteAction(c *gin.Context, actionConfig *Action, payload interface{}) (data []map[string]interface{}, err error) {
	var config *Configuration
	if v, ok := c.Get(KeyConfig); ok {
		config = v.(*Configuration)
	} else {
		return nil, ErrConfigurationMissing
	}
	contextParams, err := params(c)
	if err != nil {
		return nil, ErrParamMissing
	}

	id := valueSliceWithParamKey(*contextParams, "id")

	m, err := schema.Parse(config.Model, schemaStore, schema.NamingStrategy{})
	logrus.Print(m.PrimaryFieldDBNames[0])

	err = config.DB.Model(config.Model).Where(fmt.Sprintf("%s in ?", m.PrimaryFieldDBNames[0]), id).Delete(config.Model).Error

	return nil, err

}
