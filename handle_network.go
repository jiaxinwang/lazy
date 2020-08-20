package lazy

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/levigross/grequests"
	"github.com/sirupsen/logrus"
	"github.com/spyzhov/ajson"
)

// DefaultNetworkAction ...
func DefaultNetworkAction(c *gin.Context, actionConfig *ActionConfiguration, payload interface{}) (data []map[string]interface{}, err error) {
	payloadSet := payload.(map[string]interface{})
	method := payloadSet["_method"].(string)
	url := payloadSet["_url"].(string)
	ro := payloadSet["_request_options"].(*grequests.RequestOptions)
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
	logrus.Trace(resp.StatusCode)
	// logrus.Trace(resp.Bytes())
	logrus.Trace(resp.String())

	root, _ := ajson.Unmarshal(resp.Bytes())
	nodes, _ := root.JSONPath("$..user_id")
	for _, node := range nodes {
		logrus.Print(node.String())
	}

	return nil, nil
}
