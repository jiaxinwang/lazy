package lazy

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Configuration configs lazy values and actions
type Configuration struct {
	DB      *gorm.DB
	Table   string
	Columms string
	Model   interface{}
	Results []interface{}
	// TODO: remove it
	BeforeAction       *Action
	AfterAction        *Action
	IgnoreAssociations bool
	NeedCount          bool
	Action             []Action
	// TODO: remove it
	Before []Action
	After  []Action
}

// JSONPathMap ...
type JSONPathMap struct {
	Src     string
	Dest    string
	Default interface{}
	Remove  bool
}

// Action ...
type Action struct {
	DB                 *gorm.DB
	Table              string
	Columms            string
	Model              interface{}
	Results            []interface{}
	Params             []JSONPathMap
	ResultMaps         []JSONPathMap
	Validates          map[string]string
	IgnoreAssociations bool
	NeedCount          bool
	Payload            interface{}
	Action             func(c *gin.Context, actionConfig *Action, payload interface{}) (data []map[string]interface{}, err error)
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
