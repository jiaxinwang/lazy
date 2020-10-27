package lazy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	gm "github.com/jiaxinwang/common/gin-middleware"

	// gm "github.com/jiaxinwang/common/gin-middleware"

	"github.com/stretchr/testify/assert"

	"github.com/google/go-cmp/cmp"
)

type Ret struct {
	Count int   `json:"count"`
	Items []Dog `json:"items"`
}

func router() *gin.Engine {
	r := gin.Default()
	r.Use(gm.Trace, MiddlewareParams, MiddlewareResponse, MiddlewareDefaultResult, MiddlewareExec)
	return r
}

func buildDogMiddlewareDefaultHandlerRouter(r *gin.Engine) *gin.Engine {
	r.GET("/dogs", func(c *gin.Context) {
		config := Configuration{
			DB:        gormDB,
			Table:     "dogs",
			Columms:   "*",
			Model:     &Dog{},
			Results:   []interface{}{},
			NeedCount: true,
			Action:    []Action{{DB: gormDB, Model: &Dog{}, Action: DefaultGetAction}},
		}
		c.Set(KeyConfig, &config)
		return
	})
	r.DELETE("/dogs/:id", func(c *gin.Context) {
		config := Configuration{
			DB:        gormDB,
			Table:     "dogs",
			Columms:   "*",
			Model:     &Dog{},
			Results:   []interface{}{},
			NeedCount: true,
		}
		c.Set(KeyConfig, &config)
		return
	})
	r.POST("/dogs", func(c *gin.Context) {
		config := Configuration{
			DB:     gormDB,
			Model:  &Dog{},
			Action: []Action{{DB: gormDB, Model: &Dog{}, Action: DefaultPostAction}},
		}
		c.Set(KeyConfig, &config)
		return
	})
	r.PATCH("/dogs/:id", func(c *gin.Context) {
		config := Configuration{
			DB:     gormDB,
			Model:  &Dog{},
			Action: []Action{{DB: gormDB, Model: &Dog{}, Action: DefaultPatchAction}},
		}
		c.Set(KeyConfig, &config)
		return
	})
	// }

	return r
}

func TestDefaultHTTPActionMiddleware(t *testing.T) {
	initTestDB()
	r := router()
	g := r.Use(MiddlewareParams)
	{
		g.GET("/dogs/http", func(c *gin.Context) {
			// payload := &HTTPRequest{"https://httpbin.org/anything/1", "GET", map[string]interface{}{"k": 1}}
			payload := &HTTPRequest{"https://httpbin.org/anything/1", "GET", 1}
			config := Configuration{
				DB:        gormDB,
				Table:     "dogs",
				Columms:   "*",
				Model:     &Dog{},
				Results:   []interface{}{},
				NeedCount: true,
				Action:    []Action{{DB: gormDB, Payload: payload, Action: DefaultHTTPRequestAction}},
			}
			c.Set(KeyConfig, &config)
			return
		})
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/dogs/http", nil)
	r.ServeHTTP(w, req)
	response := Response{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 200, w.Code)
	assert.NoError(t, err)
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	return r
}

func TestGin(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestDeleteHandle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	config := Configuration{
		DB:    gormDB,
		Model: &Dog{},
	}

	configIgnoreAssociations := Configuration{
		DB:                 gormDB,
		Model:              &Dog{},
		IgnoreAssociations: true,
	}

	type args struct {
		c    *gin.Context
		id   string
		conf *Configuration
	}
	tests := []struct {
		name     string
		args     args
		wantData []map[string]interface{}
		wantErr  bool
	}{
		{"case-1", args{c: context, id: "1", conf: &configIgnoreAssociations}, nil, false},
		{"case-2", args{c: context, id: "abc", conf: &config}, nil, true},
		{"case-3", args{c: context, id: "2", conf: &config}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.c.Set(KeyConfig, tt.args.conf)
			tt.args.c.Params = []gin.Param{{
				Key:   `id`,
				Value: tt.args.id,
			}}

			gotData, err := DeleteHandle(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteHandle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(gotData, tt.wantData) {
				t.Errorf("DeleteHandle() = %v, want %v\ndiff=%v", gotData, tt.wantData, cmp.Diff(gotData, tt.wantData))
			}
		})
	}
}
