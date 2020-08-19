package lazy

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/levigross/grequests"
	"github.com/sirupsen/logrus"
)

// DefaultNetworkAction ...
func DefaultNetworkAction(c *gin.Context, actionConfig *ActionConfiguration, payload interface{}) (data []map[string]interface{}, err error) {
	payloadSet := payload.(map[string]interface{})
	method := payloadSet["_method"].(string)
	url := payloadSet["_url"].(string)
	logrus.Print(method)
	ro := payloadSet["_request_options"].(*grequests.RequestOptions)
	logrus.WithFields(
		logrus.Fields{
			"url":             url,
			"method":          method,
			"request_options": *ro,
		},
	).Trace()
	resp := &grequests.Response{}
	switch method {
	case http.MethodGet:
	case http.MethodHead:
	case http.MethodPost:
		resp, err = grequests.Post(url, ro)
	case http.MethodPut:
	case http.MethodPatch:
	case http.MethodDelete:
	case http.MethodConnect:
	case http.MethodOptions:
	case http.MethodTrace:
	}
	logrus.Trace(resp.String())
	return nil, nil
}
