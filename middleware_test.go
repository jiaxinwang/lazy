package lazy

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
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

func TestMiddlewareDefaultHandlerPOST(t *testing.T) {
	initTestDB()
	r := buildDogMiddlewareDefaultHandlerRouter(router())
	w := httptest.NewRecorder()

	body := `{"name":"test-put-dog-2","foods":[{"id":1},{"id":2}]}`

	req, _ := http.NewRequest("POST", "/dogs", bytes.NewReader([]byte(body)))

	r.ServeHTTP(w, req)
	response := Response{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		logrus.WithError(err).Debug()
	}
	assert.Equal(t, 200, w.Code)
	assert.NoError(t, err)
}
