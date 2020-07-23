package lazy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	gm "github.com/jiaxinwang/common/gin-middleware"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/google/go-cmp/cmp"
)

type Ret struct {
	Count int   `json:"count"`
	Items []Dog `json:"items"`
}

func router() *gin.Engine {
	r := gin.Default()
	r.Use(gm.Trace)
	r.Use(gm.LazyResponse)

	return r
}

func buildDogGetMiddlewareRouter(r *gin.Engine) *gin.Engine {
	r.Use(Middleware).Use(MiddlewareTransParams).GET("/dogs", func(c *gin.Context) {
		config := Configuration{
			DB:        gormDB,
			Table:     "dogs",
			Columms:   "*",
			Model:     &Dog{},
			Results:   []interface{}{},
			NeedCount: true,
		}
		c.Set("_lazy_configuration", &config)
		return
	})
	r.Use(Middleware).Use(MiddlewareTransParams).DELETE("/dog/:id", func(c *gin.Context) {
		config := Configuration{
			DB:        gormDB,
			Table:     "dogs",
			Columms:   "*",
			Model:     &Dog{},
			Results:   []interface{}{},
			NeedCount: true,
		}
		c.Set("_lazy_configuration", &config)
		return
	})

	return r
}

func buildDogGetRouter(r *gin.Engine) *gin.Engine {
	r.Use(MiddlewareTransParams).GET("/dogs", func(c *gin.Context) {
		config := Configuration{
			DB:        gormDB,
			Table:     "dogs",
			Columms:   "*",
			Model:     &Dog{},
			Results:   []interface{}{},
			NeedCount: true,
		}
		c.Set("_lazy_configuration", &config)
		if _, err := GetHandle(c); err != nil {
			c.Set("error_msg", err.Error())
			return
		}
		if v, exist := c.Get(keyResults); exist {
			c.Set("ret", map[string]interface{}{"data": v})
		}
		return
	})
	return r
}

func TestActionHandlePage(t *testing.T) {
	r := buildDogGetRouter(router())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/dogs", nil)
	q := req.URL.Query()
	q.Add(`page`, `0`)
	q.Add(`limit`, `1`)
	q.Add(`offset`, `1`)
	req.URL.RawQuery = q.Encode()

	r.ServeHTTP(w, req)
	response := Response{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 200, w.Code)
	assert.NoError(t, err)

	var ret Ret
	MapStruct(response.Data.(map[string]interface{}), &ret)
	// logrus.Printf("%+v", ret)

	assert.Equal(t, 2, ret.Count)
	assert.Equal(t, 1, len(ret.Items))

}

func TestActionHandle(t *testing.T) {
	r := buildDogGetRouter(router())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/dogs", nil)
	q := req.URL.Query()
	q.Add("id", `1`)
	q.Add("id", `2`)
	req.URL.RawQuery = q.Encode()

	r.ServeHTTP(w, req)
	response := Response{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 200, w.Code)
	assert.NoError(t, err)
}

func TestActionHandleMiddleware(t *testing.T) {
	r := buildDogGetMiddlewareRouter(router())
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/dogs", nil)
	q := req.URL.Query()
	q.Add("id", `1`)
	q.Add("id", `2`)
	req.URL.RawQuery = q.Encode()

	r.ServeHTTP(w, req)
	response := Response{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 200, w.Code)
	assert.NoError(t, err)
	var ret Ret
	MapStruct(response.Data.(map[string]interface{}), &ret)
	// logrus.Printf("%+v", ret)

	assert.Equal(t, 2, ret.Count)
	assert.Equal(t, 2, len(ret.Items))
}

func TestBeforeActionHandle(t *testing.T) {
	r := router()
	r.Use(MiddlewareTransParams).GET("/dogs", func(c *gin.Context) {
		var ret []interface{}
		config := Configuration{
			DB: gormDB,
			Before: &ActionConfiguration{
				Table:     "profiles",
				Model:     &Profile{},
				ResultMap: map[string]string{"dog_id": "id"},
				Action:    DefaultBeforeAction,
				Params:    []string{`before_dog_id`},
			},
			Table:   "dogs",
			Columms: "*",
			Model:   &Dog{},
			Results: ret,
		}
		c.Set("_lazy_configuration", &config)
		if _, err := GetHandle(c); err != nil {
			c.Set("error_msg", err.Error())
			return
		}
		c.Set("ret", map[string]interface{}{"data": map[string]interface{}{"count": len(config.Results), "items": config.Results}})
		return
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/dogs", nil)
	q := req.URL.Query()
	q.Add("before_dog_id", `1`)
	q.Add("before_dog_id", `2`)
	req.URL.RawQuery = q.Encode()

	r.ServeHTTP(w, req)
	response := Response{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 200, w.Code)
	assert.NoError(t, err)

	var ret Ret
	MapStruct(response.Data.(map[string]interface{}), &ret)
	logrus.Printf("%+v", ret)

}

// func TestAfterActionHandle(t *testing.T) {
// 	r := router()

// 	r.GET("/dogs", func(c *gin.Context) {
// 		var ret []interface{}
// 		config := Configuration{
// 			DB: gormDB,
// 			After: &ActionConfiguration{
// 				Table:     "profiles",
// 				Columms:   "dog_id",
// 				Model:     &Profile{},
// 				ResultMap: map[string]string{"dog_id": "id"},
// 				Action:    DefaultBeforeAction,
// 			},
// 			Table:   "dogs",
// 			Columms: "*",
// 			Model:   &Dog{},
// 			Results: ret,
// 		}
// 		c.Set("lazy-configuration", &config)
// 		if _, err := Handle(c); err != nil {
// 			c.Set("error_msg", err.Error())
// 			return
// 		}
// 		c.Set("ret", map[string]interface{}{"data": config.Results})
// 		return
// 	})

// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("GET", "/dogs", nil)
// 	q := req.URL.Query()
// 	q.Add("before_dog_id", `1`)
// 	q.Add("before_dog_id", `2`)
// 	req.URL.RawQuery = q.Encode()

// 	r.ServeHTTP(w, req)
// 	ret := Response{}
// 	err := json.Unmarshal(w.Body.Bytes(), &ret)
// 	assert.Equal(t, 200, w.Code)
// 	assert.NoError(t, err)
// 	logrus.Print(ret)
// }

func TestDeleteHandle1(t *testing.T) {
	r := router()
	r.Use(Middleware).Use(MiddlewareTransParams).DELETE("/dog/:id", func(c *gin.Context) {
		var ret []interface{}
		config := Configuration{
			DB:      gormDB,
			Table:   "dogs",
			Model:   &Dog{},
			Results: ret,
		}
		c.Set("_lazy_configuration", &config)
		return
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/dog", nil)

	r.ServeHTTP(w, req)
	response := Response{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 200, w.Code)
	assert.NoError(t, err)

}

func TestDeleteHandle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	type args struct {
		c  *gin.Context
		id string
	}
	tests := []struct {
		name     string
		args     args
		wantData []map[string]interface{}
		wantErr  bool
	}{
		{"case-1", args{c: context, id: "1"}, nil, false},
		{"case-2", args{c: context, id: "abc"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.c.Params = []gin.Param{gin.Param{
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
