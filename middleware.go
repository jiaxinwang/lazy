package lazy

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	// KeyParams ...
	KeyParams = `_lazy_params`
	// KeyParamsUnion ...
	KeyParamsUnion = `_lazy_params_union`
	// KeyBody ...
	KeyBody = `_lazy_body`
	// KeyConfig ...
	KeyConfig = `_lazy_configuration`
	// KeyErrorMessage ...
	KeyErrorMessage = `error_msg`
	keyResults      = `_lazy_results`
	keyCount        = `_lazy_count`
	keyData         = `_lazy_data`
)

// MiddlewareParams ...
func MiddlewareParams(c *gin.Context) {
	params := map[string][]string(c.Request.URL.Query())

	// TODO:
	if id := c.Param("id"); !strings.EqualFold(id, "") {
		params["id"] = []string{id}
	}

	c.Set(KeyParams, params)

	body := make(map[string]interface{})
	if c.Request.Body != nil {
		if b, err := ioutil.ReadAll(c.Request.Body); err != nil {
			logrus.WithError(err).Trace()
		} else {
			if json.Unmarshal(b, &body) != nil {
				logrus.WithError(err).Trace()
			}
		}
	}
	c.Set(KeyBody, body)

	union := make(map[string]interface{})
	for k, v := range params {
		union[k] = v
	}
	for k, v := range body {
		union[k] = v
	}

	c.Set(KeyParamsUnion, union)

	c.Next()
}

// MiddlewareDefaultResult ...
func MiddlewareDefaultResult(c *gin.Context) {
	defer func() {
		if v, exist := c.Get(keyResults); exist {
			c.Set(keyData, map[string]interface{}{"data": v})
		}
	}()
	c.Next()
}

// MiddlewareExec ...
func MiddlewareExec(c *gin.Context) {
	defer func() {
		var config *Configuration
		if v, ok := c.Get(KeyConfig); ok {
			config = v.(*Configuration)
		} else {
			c.Set(KeyErrorMessage, ErrConfigurationMissing)
			return
		}
		switch c.Request.Method {
		case http.MethodGet, http.MethodDelete, http.MethodPost, http.MethodPatch, http.MethodPut:
			for _, v := range config.Action {
				if _, err := v.Action(c, &v, v.Payload); err != nil {
					c.Set(KeyErrorMessage, err.Error())
				}
			}
		}
		return
	}()
	c.Next()
}

// MiddlewareResponse ...
func MiddlewareResponse(c *gin.Context) {
	defer func() {
		ret := make(map[string]interface{})
		if v, ok := c.Get(keyData); ok {
			ret = v.(map[string]interface{})
		}
		if _, ok := ret["data"]; !ok {
			ret["data"] = nil
		}
		if _, ok := ret["error_no"]; !ok {
			if _, ok := ret[KeyErrorMessage]; ok {
				ret["error_no"] = 400
			} else {
				ret["error_no"] = 0
			}
		}
		if _, ok := ret[KeyErrorMessage]; !ok {
			ret[KeyErrorMessage] = ``
		}

		ret["request_id"] = c.MustGet("requestID")
		c.JSON(200, ret)
	}()
	c.Next()
}
