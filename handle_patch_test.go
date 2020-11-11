package lazy

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultPatchAction(t *testing.T) {
	initTestDB()
	r := defaultDogRouter(router())

	jsonParams := []string{
		// `{"name":"test-patch-dog-1"}`,
		// `{"name":"test-patch-dog-2","foods":[{"id":1},{"id":2}]}`,
		// `{"name":"test-patch-dog-3","toys":[{"name":"new-toy-1"},{"name":"new-toy-1"}]}`,
		`{"name":"test-patch-dog-3","owner":[{"id":"1"}}`,
	}

	for _, jsonParam := range jsonParams {
		w := httptest.NewRecorder()
		contentBuffer := bytes.NewBuffer([]byte(jsonParam))
		req, _ := http.NewRequest(http.MethodPatch, "/dogs/1", contentBuffer)

		r.ServeHTTP(w, req)
		response := Response{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, 200, w.Code)
		assert.NoError(t, err)
	}
}
