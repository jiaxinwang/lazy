package lazy

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Configuration configs lazy values and actions
type Configuration struct {
	DB        *gorm.DB
	Model     interface{}
	NeedCount bool
	Action    []Action
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
	DB         *gorm.DB
	Model      interface{}
	Params     []JSONPathMap
	ResultMaps []JSONPathMap
	Validates  map[string]string
	NeedCount  bool
	Payload    interface{}
	Action     func(c *gin.Context, actionConfig *Action, payload interface{}) (data []map[string]interface{}, err error)
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
