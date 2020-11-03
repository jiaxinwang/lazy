package lazy

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

func genContent() *gin.Context {
	w := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(w)
	return context
}

func TestDefaultPostAction(t *testing.T) {
	initTestDB()
	r := defaultDogRouter(router())
	jsonParams := []string{
		`{"name":"test-post-dog-1"}`,
		`{"name":"test-post-dog-2","foods":[{"id":1},{"id":2}]}`,
	}

	for _, jsonParam := range jsonParams {
		w := httptest.NewRecorder()

		contentBuffer := bytes.NewBuffer([]byte(jsonParam))
		req, _ := http.NewRequest("POST", "/dogs", contentBuffer)

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
