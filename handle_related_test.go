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

func TestDefaultRelatedHasManyPostAction(t *testing.T) {
	initTestDB()
	r := defaultDogRouter(router())
	jsonParams := []string{
		`{"toys":[{"name":"new-toy-a"},{"name":"new-toy-b"}]}`,
	}

	for _, jsonParam := range jsonParams {
		w := httptest.NewRecorder()
		contentBuffer := bytes.NewBuffer([]byte(jsonParam))
		req, _ := http.NewRequest("POST", "/dogs/1/toys", contentBuffer)

		r.ServeHTTP(w, req)
		response := Response{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, 200, w.Code)
		assert.NoError(t, err)

		dog := Dog{}
		json.Unmarshal([]byte(jsonParam), &dog)

		dbDog := Dog{}
		if err = gormDB.Where("id = ?", 1).Preload("Toys").Find(&dbDog).Error; err != nil {
			t.Errorf("db find dog = %v", err)
			return
		}

		ignoreDog := []string{"ID", "CreatedAt", "UpdatedAt", "Foods", "Name", "BreedID", "OwnerID"}
		ignoreToy := []string{"ID", "CreatedAt", "UpdatedAt", "DogID"}

		if !cmp.Equal(
			dog,
			dbDog,
			cmpopts.IgnoreFields(Dog{}, ignoreDog...),
			cmpopts.IgnoreFields(Toy{}, ignoreToy...),
		) {
			t.Errorf(
				"dog() = %+v, want %+v\ndiff=%+v",
				dog, dbDog,
				cmp.Diff(
					dog,
					dbDog,
					cmpopts.IgnoreFields(Dog{}, ignoreDog...),
					cmpopts.IgnoreFields(Toy{}, ignoreToy...),
				),
			)
		}
	}
}

func TestDefaultRelatedMany2ManyPostAction(t *testing.T) {
	initTestDB()
	r := defaultDogRouter(router())
	jsonParams := []string{
		`{"foods":[{"id":1},{"name":"Diamond"}]}`,
	}

	for _, jsonParam := range jsonParams {
		w := httptest.NewRecorder()
		contentBuffer := bytes.NewBuffer([]byte(jsonParam))
		req, _ := http.NewRequest("POST", "/dogs/1/foods", contentBuffer)

		r.ServeHTTP(w, req)
		response := Response{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, 200, w.Code)
		assert.NoError(t, err)

		dog := Dog{}
		json.Unmarshal([]byte(jsonParam), &dog)
		dog.Foods = []Food{{ID: 1, Brand: "Pedigree"}}

		dbDog := Dog{}
		if err = gormDB.Where("id = ?", 1).Preload("Foods").Find(&dbDog).Error; err != nil {
			t.Errorf("db find dog = %v", err)
			return
		}

		ignoreDog := []string{"ID", "CreatedAt", "UpdatedAt", "Toys", "Name", "BreedID", "OwnerID"}
		ignoreFood := []string{"ID", "CreatedAt", "UpdatedAt"}

		if !cmp.Equal(
			dog,
			dbDog,
			cmpopts.IgnoreFields(Dog{}, ignoreDog...),
			cmpopts.IgnoreFields(Food{}, ignoreFood...),
		) {
			t.Errorf(
				"dog() = %+v, want %+v\ndiff=%+v",
				dog, dbDog,
				cmp.Diff(
					dog,
					dbDog,
					cmpopts.IgnoreFields(Dog{}, ignoreDog...),
					cmpopts.IgnoreFields(Food{}, ignoreFood...),
				),
			)
		}
	}
}
