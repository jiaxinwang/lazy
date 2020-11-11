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
	r := defaultDogRouter(router())
	w := httptest.NewRecorder()

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

func TestDefaultGetActionParams(t *testing.T) {
	initTestDB()
	r := defaultDogRouter(router())
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/dogs", nil)

	q := req.URL.Query()
	q.Add("name", `Max`)
	req.URL.RawQuery = q.Encode()

	var dog1 Dog
	gormDB.Where("name = Max").Preload("Toys").Preload("Foods").Preload("Owner").Find(&dog1)
	logrus.Printf("%#v", dog1)

	r.ServeHTTP(w, req)
	response := Response{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 200, w.Code)
	assert.NoError(t, err)
	var ret Ret
	MapStruct(response.Data.(map[string]interface{}), &ret)
	logrus.WithField("ret", fmt.Sprintf("%+v", ret)).Info()

}

func TestDefaultGetActionParamsHasMany(t *testing.T) {
	initTestDB()
	r := defaultDogRouter(router())
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/dogs", nil)

	q := req.URL.Query()
	q.Add("toys", `2`)
	req.URL.RawQuery = q.Encode()

	r.ServeHTTP(w, req)
	response := Response{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 200, w.Code)
	assert.NoError(t, err)
	var ret Ret
	MapStruct(response.Data.(map[string]interface{}), &ret)

	var results []map[string]interface{}
	var dbRet Dog

	gormDB.Model(&Dog{}).Joins("Toys").Where("Toys__id = 2").Find(&results)

	assert.Equal(t, len(results), 1)
	MapStruct(results[0], &dbRet)

	assert.Equal(t, dbRet.ID, ret.Items[0].ID)
}

func TestDefaultGetActionParamsMany2Many(t *testing.T) {
	initTestDB()
	r := defaultDogRouter(router())
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/dogs", nil)

	q := req.URL.Query()
	q.Add("foods", `2`)
	req.URL.RawQuery = q.Encode()

	r.ServeHTTP(w, req)
	response := Response{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 200, w.Code)
	assert.NoError(t, err)
	var ret Ret
	MapStruct(response.Data.(map[string]interface{}), &ret)

	// logrus.Printf("%+v", ret.Items)

	for _, v := range ret.Items {
		logrus.Print(v.ID)
		for _, vFood := range v.Foods {
			if vFood.ID == 2 {
				logrus.Print("-->", v.ID)
			}
		}
	}

	// var results []map[string]interface{}
	// var dbRet Dog
	// gormDB.Model(&Dog{}).Joins("Foods").Where("Foods__id = 2").Find(&results)
	// logrus.Print(results)

	// assert.Equal(t, len(results), 1)
	// MapStruct(results[0], &dbRet)

	// assert.Equal(t, dbRet.ID, ret.Items[0].ID)
}
