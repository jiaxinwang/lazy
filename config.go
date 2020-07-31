package lazy

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Configuration configs lazy values and actions
type Configuration struct {
	DB                 *gorm.DB
	Table              string
	Columms            string
	Model              interface{}
	Results            []interface{}
	BeforeAction       *ActionConfiguration
	AfterAction        *ActionConfiguration
	IgnoreAssociations bool
	NeedCount          bool
	Action             []*ActionConfiguration
	Before             []*ActionConfiguration
	After              []*ActionConfiguration
}

// ActionConfiguration configs action, before-action, after-action values and actions
type ActionConfiguration struct {
	DB                 *gorm.DB
	Table              string
	Columms            string
	Model              interface{}
	Results            []interface{}
	Params             []string
	ResultMap          map[string]string
	IgnoreAssociations bool
	NeedCount          bool
	Action             func(c *gin.Context, actionConfig *ActionConfiguration, payload interface{}) (result interface{}, reduce map[string][]string, err error)
}
