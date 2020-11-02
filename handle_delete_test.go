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
	// q.Add("id", `1`)
	// q.Add("id", `2`)
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

	logrus.Print(dog)

	// err = gormDB.Where("id = 2").Take(&dog).Error
	// assert.Equal(t, gorm.ErrRecordNotFound, err)
}
