package lazy

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

func TestDefaultPutAction(t *testing.T) {
	initTestDB()
	r := defaultDogRouter(router())
	jsonParams := []string{
		`{"name":"test-put-dog-1"}`,
		`{"name":"test-put-dog-2","foods":[{"id":1},{"id":2}]}`,
	}

	for _, jsonParam := range jsonParams {
		w := httptest.NewRecorder()
		contentBuffer := bytes.NewBuffer([]byte(jsonParam))
		req, _ := http.NewRequest("PUT", "/dogs/1", contentBuffer)

		// q := req.URL.Query()
		// q.Add("id", `1`)
		// req.URL.RawQuery = q.Encode()

		r.ServeHTTP(w, req)
		response := Response{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, 200, w.Code)
		assert.NoError(t, err)

		dog := Dog{}
		json.Unmarshal([]byte(jsonParam), &dog)
		if dog.Foods == nil {
			dog.Foods = []Food{}
		}
		dbDog := Dog{}
		if err = gormDB.Where("name = ?", dog.Name).Preload("Foods").Find(&dbDog).Error; err != nil {
			t.Errorf("db find dog = %v", err)
			return
		}

		if !cmp.Equal(
			dog,
			dbDog,
			cmpopts.IgnoreFields(Dog{}, "ID", "CreatedAt", "UpdatedAt"),
			cmpopts.IgnoreFields(Food{}, "ID", "CreatedAt", "UpdatedAt", "Brand"),
		) {
			t.Errorf(
				"dog() = %v, want %v\ndiff=%v",
				dog, dbDog,
				cmp.Diff(
					dog,
					dbDog,
					cmpopts.IgnoreFields(Dog{}, "ID", "CreatedAt", "UpdatedAt"),
					cmpopts.IgnoreFields(Food{}, "ID", "CreatedAt", "UpdatedAt", "Brand"),
				),
			)
		}
	}
}
