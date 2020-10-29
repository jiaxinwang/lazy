package lazy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestDefaultDeleteAction(t *testing.T) {
	initTestDB()
	r := buildDogMiddlewareDefaultHandlerRouter(router())
	w := httptest.NewRecorder()

	// GET
	req, _ := http.NewRequest("DELETE", "/dogs/:id", nil)

	q := req.URL.Query()
	q.Add("id", `1`)
	q.Add("id", `2`)
	req.URL.RawQuery = q.Encode()

	r.ServeHTTP(w, req)
	response := Response{}
	logrus.WithField("string", w.Body.String()).Trace()
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 200, w.Code)
	assert.NoError(t, err)

	var dog Dog
	err = gormDB.Where("id = 1").Take(&dog).Error
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	err = gormDB.Where("id = 2").Take(&dog).Error
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
