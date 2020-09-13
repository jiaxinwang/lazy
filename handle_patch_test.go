package lazy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// initTestDB()
// r := buildDogMiddlewareRouter(router())
// w := httptest.NewRecorder()
// req, _ := http.NewRequest("GET", "/dogs", nil)
// q := req.URL.Query()
// q.Add("id", `1`)
// q.Add("id", `2`)
// req.URL.RawQuery = q.Encode()

// r.ServeHTTP(w, req)
// response := Response{}
// err := json.Unmarshal(w.Body.Bytes(), &response)
// assert.Equal(t, 200, w.Code)
// assert.NoError(t, err)
// var ret Ret
// MapStruct(response.Data.(map[string]interface{}), &ret)

func buildPatchDog(r *gin.Engine) *gin.Engine {
	g := r.Use(MiddlewareParams).Use(MiddlewareExec)
	{
		g.PATCH("/dogs/:id", func(c *gin.Context) {
			// config := Configuration{
			// 	DB:     gormDB,
			// 	Model:  &Dog{},
			// 	Action: []ActionConfiguration{{DB: gormDB, Model: &Dog{}, Action: DefaultPatchAction}},
			// }
			// c.Set(KeyConfig, &config)
			return
		})
	}
	return r
}

func TestDefaultPatchAction(t *testing.T) {
	initTestDB()
	r := buildPatchDog(router())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/dogs/:id", nil)
	r.ServeHTTP(w, req)
	response := Response{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 200, w.Code)
	assert.NoError(t, err)

	// jsonParams := []string{
	// 	`{"name":"test-patch-dog-1"}`,
	// }

	// for _, jsonParam := range jsonParams {
	// 	w := httptest.NewRecorder()

	// 	contentBuffer := bytes.NewBuffer([]byte(jsonParam))
	// 	req, _ := http.NewRequest(http.MethodPatch, "/1", contentBuffer)

	// 	r.ServeHTTP(w, req)
	// 	response := Response{}
	// 	err := json.Unmarshal(w.Body.Bytes(), &response)
	// 	assert.Equal(t, 200, w.Code)
	// 	assert.NoError(t, err)

	// 	// dog := Dog{}
	// 	// json.Unmarshal([]byte(jsonParam), &dog)
	// 	// if dog.Foods == nil {
	// 	// 	dog.Foods = []Food{}
	// 	// }
	// 	// dbDog := Dog{}
	// 	// if err = gormDB.Where("name = ?", dog.Name).Find(&dbDog).Error; err != nil {
	// 	// 	t.Errorf("db find dog = %v", err)
	// 	// 	return
	// 	// }
	// 	// gormDB.Model(&dbDog).Related(&dbDog.Foods, "Foods")

	// 	// if !cmp.Equal(
	// 	// 	dog,
	// 	// 	dbDog,
	// 	// 	cmpopts.IgnoreFields(Dog{}, "ID", "CreatedAt", "UpdatedAt"),
	// 	// 	cmpopts.IgnoreFields(Food{}, "ID", "CreatedAt", "UpdatedAt", "Brand"),
	// 	// ) {
	// 	// 	t.Errorf(
	// 	// 		"dog() = %v, want %v\ndiff=%v",
	// 	// 		dog, dbDog,
	// 	// 		cmp.Diff(
	// 	// 			dog,
	// 	// 			dbDog,
	// 	// 			cmpopts.IgnoreFields(Dog{}, "ID", "CreatedAt", "UpdatedAt"),
	// 	// 			cmpopts.IgnoreFields(Food{}, "ID", "CreatedAt", "UpdatedAt", "Brand"),
	// 	// 		),
	// 	// 	)
	// 	// }
	// }
}
