package lazy

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	// KeyParams ...
	KeyParams = `_lazy_params`
	// KeyBody ...
	KeyBody    = `_lazy_body`
	keyResults = `_lazy_results`
	keyCount   = `_lazy_count`
	keyData    = `_lazy_data`
	keyConfig  = `_lazy_configuration`
)

// MiddlewareParams ...
func MiddlewareParams(c *gin.Context) {
	params := Params(c.Request.URL.Query())
	c.Set(KeyParams, params)

	body := make(map[string]interface{})

	if b, err := ioutil.ReadAll(c.Request.Body); err != nil {
		logrus.WithError(err).Trace()
	} else {
		if json.Unmarshal(b, &body) != nil {
			logrus.WithError(err).Trace()
		} else {
			c.Set(KeyBody, body)
		}
	}

	c.Next()
}

// Middleware ...
func Middleware(c *gin.Context) {
	defer func() {
		var config *Configuration
		if v, ok := c.Get(keyConfig); ok {
			config = v.(*Configuration)
		} else {
			c.Set("error_msg", ErrNoConfiguration)
			return
		}
		switch c.Request.Method {
		case http.MethodGet:
			for _, v := range config.Action {
				if _, err := v.Action(c, &v, v.Payload); err != nil {
					c.Set("error_msg", err.Error())
				}
			}
			if data, exist := c.Get(keyResults); exist {
				c.Set("ret", map[string]interface{}{"data": data})
			}
		case http.MethodDelete:
			if _, err := DeleteHandle(c); err != nil {
				c.Set("error_msg", err.Error())
				return
			}
			if data, exist := c.Get(keyResults); exist {
				c.Set("ret", map[string]interface{}{"data": data})
			}
		case http.MethodPost:
			for _, v := range config.Action {
				if _, err := v.Action(c, &v, v.Payload); err != nil {
					c.Set("error_msg", err.Error())
				}
			}

			if data, exist := c.Get(keyResults); exist {
				c.Set("ret", map[string]interface{}{"data": data})
				return
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
		if v, ok := c.Get("ret"); ok {
			ret = v.(map[string]interface{})
		}
		if _, ok := ret["data"]; !ok {
			ret["data"] = nil
		}
		if _, ok := ret["error_no"]; !ok {
			if _, ok := ret["error_msg"]; ok {
				ret["error_no"] = 400
			} else {
				ret["error_no"] = 0
			}
		}
		if _, ok := ret["error_msg"]; !ok {
			ret["error_msg"] = ``
		}

		ret["request_id"] = c.MustGet("requestID")
		c.JSON(200, ret)
	}()
	c.Next()
}
