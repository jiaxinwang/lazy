package lazy

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	keyParams  = `_lazy_params`
	keyResults = `_lazy_results`
	keyCount   = `_lazy_count`
	keyData    = `_lazy_data`
	keyConfig  = `_lazy_configuration`
)

// MiddlewareTransParams trans params into content
func MiddlewareTransParams(c *gin.Context) {
	params := Params(c.Request.URL.Query())
	c.Set(keyParams, params)
	c.Next()
}

// Middleware run the query
func Middleware(c *gin.Context) {
	defer func() {
		switch c.Request.Method {
		case http.MethodGet:
			if _, err := GetHandle(c); err != nil {
				c.Set("error_msg", err.Error())
				return
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
			if _, err := PostHandle(c); err != nil {
				c.Set("error_msg", err.Error())
				return
			}
			if data, exist := c.Get(keyResults); exist {
				c.Set("ret", map[string]interface{}{"data": data})
			}
		}
		return
	}()
	c.Next()
}
