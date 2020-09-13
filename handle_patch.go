package lazy

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// DefaultPatchAction ...
func DefaultPatchAction(c *gin.Context, actionConfig *ActionConfiguration, payload interface{}) (data []map[string]interface{}, err error) {
	var config *Configuration
	if config, err = configuration(c); err != nil {
		return nil, errors.Wrap(err, ErrConfigurationMissing.Error())
	}

	var p *Params
	if p, err = params(c); err != nil {
		return nil, errors.Wrap(err, ErrConfigurationMissing.Error())
	}

	eq, _, _, _, _ := URLValues(config.Model, *p)
	logrus.Print(eq)

	c.Set(keyResults, map[string]interface{}{"code": 0})

	return
}
