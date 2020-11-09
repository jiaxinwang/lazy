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
	r := defaultDogRouter(router())
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("DELETE", "/dogs/6", nil)

	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()

	r.ServeHTTP(w, req)
	response := Response{}
	logrus.WithField("string", w.Body.String()).Trace()
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 200, w.Code)
	assert.NoError(t, err)

	var dog Dog
	err = gormDB.Where("id = 6").Take(&dog).Error
	assert.Equal(t, gorm.ErrRecordNotFound, err)

}

func TestDefaultDeleteActionMany2Many(t *testing.T) {
	initTestDB()
	r := defaultDogRouter(router())
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("DELETE", "/foods/1", nil)

	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()

	r.ServeHTTP(w, req)
	response := Response{}
	logrus.WithField("string", w.Body.String()).Trace()
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 200, w.Code)
	assert.NoError(t, err)

	var food Food
	err = gormDB.Where("id = 1").Take(&food).Error
	assert.Equal(t, gorm.ErrRecordNotFound, err)

}

func TestDefaultDeleteActionHasMany(t *testing.T) {
	initTestDB()
	r := defaultDogRouter(router())
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("DELETE", "/toys/1", nil)

	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()

	r.ServeHTTP(w, req)
	response := Response{}
	logrus.WithField("string", w.Body.String()).Trace()
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 200, w.Code)
	assert.NoError(t, err)

	var toy Toy
	err = gormDB.Where("id = 1").Take(&toy).Error
	assert.Equal(t, gorm.ErrRecordNotFound, err)

}
