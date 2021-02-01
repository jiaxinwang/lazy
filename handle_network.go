package lazy

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/levigross/grequests"
	"github.com/sirupsen/logrus"

	"github.com/jiaxinwang/common/logger"
)

// DefaultNetworkAction ...
func DefaultNetworkAction(c *gin.Context, actionConfig *Action, payload interface{}) (data []map[string]interface{}, err error) {
	payloadSet := payload.(map[string]interface{})
	method := payloadSet["_method"].(string)
	url := payloadSet["_url"].(string)
	ro := payloadSet["_request_options"].(*grequests.RequestOptions)

	body, _ := c.Get(KeyBody)
	src := body.(map[string]interface{})
	dest := make(map[string]interface{})
	for _, v := range actionConfig.Params {
		if src, dest, err = SetJSONSingle(src, dest, v); err != nil {
			logrus.WithError(err).Trace()
		}
	}

	bodyByte, _ := json.Marshal(dest)
	ro.JSON = bodyByte

	logrus.WithField("RequestOption.JSON", logger.Pretty(dest)).Trace()
	resp := &grequests.Response{}

	switch method {
	case http.MethodGet:
	case http.MethodHead:
	case http.MethodPost:
		if resp, err = grequests.Post(url, ro); err != nil {
			logrus.WithError(err).Error()
			c.Set(KeyErrorMessage, err.Error())
			return
		}
	case http.MethodPut:
	case http.MethodPatch:
	case http.MethodDelete:
	case http.MethodConnect:
	case http.MethodOptions:
	case http.MethodTrace:
	}

	respStruct := make(map[string]interface{})
	json.UnmarshalFromString(resp.String(), &respStruct)
	ret := make(map[string]interface{})

	if actionConfig.ResultMaps != nil {
		for _, v := range actionConfig.ResultMaps {
			if respStruct, ret, err = SetJSONSingle(respStruct, ret, v); err != nil {
				logrus.WithError(err).Trace()
			}
		}
	}

	c.Set(KeyResults, ret)
	logrus.WithField(KeyResults, logger.Pretty(ret)).Trace()
	return []map[string]interface{}{ret}, nil
}
