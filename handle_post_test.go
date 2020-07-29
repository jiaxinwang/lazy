package lazy

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
)

func TestPostHandle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(w)

	// config := Configuration{
	// 	DB:    gormDB,
	// 	Model: &Dog{},
	// }

	configIgnoreAssociations := Configuration{
		DB:                 gormDB,
		Model:              &Dog{},
		IgnoreAssociations: true,
	}

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
		{"case-simple", args{c: context, json: `{"name":"test-dog"}`, conf: &configIgnoreAssociations}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.c.Set(keyConfig, tt.args.conf)
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
		})
	}
}
