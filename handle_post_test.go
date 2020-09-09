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
	r := buildDogMiddlewareRouter(router())
	jsonParams := []string{
		`{"name":"test-put-dog-1"}`,
		`{"name":"test-put-dog-2","foods":[{"id":1},{"id":2}]}`,
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
		if err = gormDB.Where("name = ?", dog.Name).Find(&dbDog).Error; err != nil {
			t.Errorf("db find dog = %v", err)
			return
		}
		gormDB.Model(&dbDog).Related(&dbDog.Foods, "Foods")

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

func TestPostAction(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type args struct {
		c    *gin.Context
		json string
		conf *Configuration
	}
	tests := []struct {
		name     string
		args     args
		wantData []map[string]interface{}
		wantErr  bool
	}{
		{"case-simple", args{c: genContent(), json: `{"name":"test-put-dog-1"}`, conf: &Configuration{DB: gormDB, Model: &Dog{}, IgnoreAssociations: true}}, nil, false},
		{"case-nested", args{c: genContent(), json: `{"name":"test-put-dog-2","foods":[{"id":1},{"id":2}]}`, conf: &Configuration{DB: gormDB, Model: &Dog{}, IgnoreAssociations: true}}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.c.Set(KeyConfig, tt.args.conf)
			contentBuffer := bytes.NewBuffer([]byte(tt.args.json))
			tt.args.c.Request, _ = http.NewRequest("POST", "/dogs", contentBuffer)

			gotData, err := PostHandle(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostHandle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !cmp.Equal(gotData, tt.wantData) {
				t.Errorf("PostHandle() = %v, want %v\ndiff=%v", gotData, tt.wantData, cmp.Diff(gotData, tt.wantData))
			}
			dog := Dog{}
			json.Unmarshal([]byte(tt.args.json), &dog)
			if dog.Foods == nil {
				dog.Foods = []Food{}
			}
			dbDog := Dog{}
			if err = gormDB.Where("name = ?", dog.Name).Find(&dbDog).Error; err != nil {
				t.Errorf("db find dog = %v", err)
				return
			}
			gormDB.Model(&dbDog).Related(&dbDog.Foods, "Foods")

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

		})
	}
}
