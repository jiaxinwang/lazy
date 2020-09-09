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
	Action             []ActionConfiguration
	Before             []ActionConfiguration
	After              []ActionConfiguration
}

// JSONPathMap ...
type JSONPathMap struct {
	Src     string
	Dest    string
	Default interface{}
	Remove  bool
}

// ActionConfiguration configs action, before-action, after-action values and actions
type ActionConfiguration struct {
	DB                 *gorm.DB
	Table              string
	Columms            string
	Model              interface{}
	Results            []interface{}
	Params             []JSONPathMap
	ResultMaps         []JSONPathMap
	IgnoreAssociations bool
	NeedCount          bool
	Payload            interface{}
	Action             func(c *gin.Context, actionConfig *ActionConfiguration, payload interface{}) (data []map[string]interface{}, err error)
}

// HTTPRequest ...
type HTTPRequest struct {
	RequestURL    string
	RequestMethod string
	RequestBody   interface{}
}

func configuration(c *gin.Context) (*Configuration, error) {
	if v, ok := c.Get(KeyConfig); ok {
		return v.(*Configuration), nil
	}
	return nil, ErrConfigurationMissing
}
