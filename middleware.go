package lazy

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
				if _, err := v.Action(c, &v, nil); err != nil {
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
