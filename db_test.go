package lazy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_createModel(t *testing.T) {
	initTestDB()
	var dog Dog
	assert.NoError(t, createModel(gormDB, &Owner{Name: "has-one-owner", Dog: Dog{Name: "has-one-dog"}}))

	// has one
	var owner Owner
	assert.NoError(t, gormDB.Preload("Dog").Where("name = ?", "has-one-owner").Find(&owner).Error)
	// assert.NoError(t, gormDB.Model(&owner).Related(&owner.Dog).Error)

	assert.Equal(t, "has-one-owner", owner.Name)
	assert.Equal(t, "has-one-dog", owner.Dog.Name)
	assert.Equal(t, owner.ID, owner.Dog.OwnerID)

	// has many
	assert.NoError(t, createModel(gormDB, &Dog{Name: "has-many-toys-dog", Toys: []Toy{{Name: "toy-1"}, {Name: "toy-2"}}}))

	assert.NoError(t, gormDB.Preload("Toys").Where("name = ?", "has-many-toys-dog").Find(&dog).Error)
	// assert.NoError(t, gormDB.Model(&dog).Related(&dog.Toys).Error)

	assert.Equal(t, "has-many-toys-dog", dog.Name)
	assert.Equal(t, len(dog.Toys), 2)
	assert.Equal(t, dog.Toys[0].Name, "toy-1")
	assert.Equal(t, dog.Toys[1].Name, "toy-2")

	//many tp many
	assert.NoError(t, createModel(gormDB, &Dog{Name: "many-to-many-dog-1", Foods: []Food{{ID: 1}, {ID: 2}}}))
	assert.NoError(t, createModel(gormDB, &Dog{Name: "many-to-many-dog-2", Foods: []Food{{ID: 3}, {ID: 1}}}))

	dog = Dog{}
	assert.NoError(t, gormDB.Preload("Foods").Where("name = ?", "many-to-many-dog-1").Find(&dog).Error)
	// assert.NoError(t, gormDB.Model(&dog).Related(&dog.Foods, "Foods").Error)
	assert.Equal(t, 2, len(dog.Foods))
	assert.Equal(t, uint(1), dog.Foods[0].ID)
	assert.Equal(t, `Pedigree`, dog.Foods[0].Brand)
	assert.Equal(t, uint(2), dog.Foods[1].ID)
	assert.Equal(t, `Purina`, dog.Foods[1].Brand)

	dog = Dog{}
	assert.NoError(t, gormDB.Preload("Foods").Where("name = ?", "many-to-many-dog-2").Find(&dog).Error)
	// assert.NoError(t, gormDB.Model(&dog).Related(&dog.Foods, "Foods").Error)
	assert.Equal(t, 2, len(dog.Foods))
	assert.Equal(t, uint(1), dog.Foods[0].ID)
	assert.Equal(t, `Pedigree`, dog.Foods[0].Brand)
	assert.Equal(t, uint(3), dog.Foods[1].ID)
	assert.Equal(t, `Diamond`, dog.Foods[1].Brand)

}

func Test_queryAssociated(t *testing.T) {
	initTestDB()
	gotRet := queryAssociated(gormDB, "dogs", "id", 1)
	assert.Equal(t, len(gotRet), 1)
}
