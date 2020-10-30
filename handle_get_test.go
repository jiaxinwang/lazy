package lazy

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestDefaultGetAction(t *testing.T) {
	initTestDB()
	r := buildDogMiddlewareDefaultHandlerRouter(router())
	w := httptest.NewRecorder()

	// GET
	req, _ := http.NewRequest("GET", "/dogs", nil)

	q := req.URL.Query()
	// q.Add("id", `1`)
	// q.Add("id", `2`)
	q.Add("id", `9`)
	req.URL.RawQuery = q.Encode()

	var dog1 Dog
	gormDB.Where("id = 1").Preload("Toys").Preload("Foods").Preload("Owner").Find(&dog1)
	logrus.Printf("%#v", dog1)

	r.ServeHTTP(w, req)
	response := Response{}
	logrus.WithField("string", w.Body.String()).Debug()
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 200, w.Code)
	assert.NoError(t, err)
	var ret Ret
	MapStruct(response.Data.(map[string]interface{}), &ret)
	logrus.WithField("ret", fmt.Sprintf("%+v", ret)).Info()

	// assert.Equal(t, 2, ret.Count)
	// assert.Equal(t, 2, len(ret.Items))

	// assert.Equal(t, ret.Items[0].ID, uint(1))
	// assert.Equal(t, ret.Items[1].ID, uint(2))

	// assert.Equal(t, len(ret.Items[0].Toys), 2)
	// assert.Equal(t, len(ret.Items[1].Toys), 2)

	// gormDB.Model(&dog1).Related(&(dog1.Toys))

	// logrus.Printf("%+v", ret)

	// assert.Equal(t, ret.Items[0].Name, dog1.Name)

}
