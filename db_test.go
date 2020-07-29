package lazy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_createModel(t *testing.T) {
	err := createModel(gormDB, &Owner{Name: "has-one-owner", Dog: Dog{Name: "has-one-dog"}})
	assert.NoError(t, err)

	var owner Owner
	// var dog Dog
	assert.NoError(t, gormDB.Where("name = ?", "has-one-owner").Find(&owner).Error)
	assert.NoError(t, gormDB.Model(&owner).Related(&owner.Dog).Error)

	assert.Equal(t, owner.Name, "has-one-owner")
	assert.Equal(t, owner.Dog.Name, "has-one-dog")
	assert.Equal(t, owner.ID, owner.Dog.OwnerID)
}
