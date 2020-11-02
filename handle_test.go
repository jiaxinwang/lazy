package lazy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	gm "github.com/jiaxinwang/common/gin-middleware"

	"github.com/stretchr/testify/assert"
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
			DB:      gormDB,
			Table:   "dogs",
			Columms: "*",
			Model:   &Dog{},
			Action:  []Action{{DB: gormDB, Model: &Dog{}, Action: DefaultDeleteAction}},
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
	r.PUT("/dogs/:id", func(c *gin.Context) {
		config := Configuration{
			DB:     gormDB,
			Table:  "dogs",
			Model:  &Dog{},
			Action: []Action{{DB: gormDB, Model: &Dog{}, Action: DefaultPutAction}},
		}
		c.Set(KeyConfig, &config)
		return
	})

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
