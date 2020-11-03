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
		`{"name":"patch-dog-name-1"}`,
		`{"name":"test-put-dog-2","foods":[{"id":1},{"id":2}]}`,
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
