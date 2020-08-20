package lazy

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func TestMiddlewareParams(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	gin.SetMode(gin.TestMode)

	type args struct {
		c    *gin.Context
		url  string
		json string
	}
	tests := []struct {
		name string
		args args
	}{
		{"case-1", args{genContent(), `/?name=1`, `{"name":"test-put-dog-1"}`}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contentBuffer := bytes.NewBuffer([]byte(tt.args.json))
			tt.args.c.Request, _ = http.NewRequest("POST", tt.args.url, contentBuffer)
			MiddlewareParams(tt.args.c)
		})
	}
}
