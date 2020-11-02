package lazy

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// DefaultPatchAction ...
func DefaultPatchAction(c *gin.Context, actionConfig *Action, payload interface{}) (data []map[string]interface{}, err error) {
	id := c.Param("id")
	logrus.WithField("id", id).Info()
	// c.Params
	// _, _, bodyParams := ContentParams(c)
	// var config *Configuration
	// if v, ok := c.Get(KeyConfig); ok {
	// 	config = v.(*Configuration)
	// } else {
	// 	return nil, ErrConfigurationMissing
	// }

	// s, err := json.MarshalToString(bodyParams)
	// if err != nil {
	// 	logrus.WithError(err).Error()
	// 	// TODO: error
	// 	return nil, err
	// }
	// err = json.UnmarshalFromString(s, &config.Model)
	// if err != nil {
	// 	logrus.WithError(err).Error()
	// 	// TODO: error
	// 	return nil, err
	// }

	// err = createModel(config.DB, config.Model)
	// data = make([]map[string]interface{}, 1)
	// data[0] = make(map[string]interface{})
	// data[0][keyData] = clone(config.Model)
	// c.Set(keyResults, data)
	// return data, err
	return nil, nil
}
