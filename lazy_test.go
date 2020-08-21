package lazy

import (
	"os"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type Response struct {
	Data     interface{} `json:"data"`
	ErrorMsg string      `json:KeyErrorMessage`
	ErrorNo  int         `json:"error_no"`
}

// Food many to many dog
type Food struct {
	ID        uint      `gorm:"primarykey" json:"id" lazy:"id" mapstructure:"id"`
	CreatedAt time.Time `json:"created_at" lazy:"created_at" mapstructure:"created_at"`
	UpdatedAt time.Time `json:"updated_at" lazy:"updated_at" mapstructure:"updated_at"`
	Brand     string    `json:"brand" lazy:"brand" mapstructure:"brand"`
}

// Owner has one dog
type Owner struct {
	ID        uint      `gorm:"primarykey" json:"id" lazy:"id" mapstructure:"id"`
	CreatedAt time.Time `json:"created_at" lazy:"created_at" mapstructure:"created_at"`
	UpdatedAt time.Time `json:"updated_at" lazy:"updated_at" mapstructure:"updated_at"`
	Name      string    `json:"name" lazy:"name" mapstructure:"name"`
	Dog       Dog       `json:"dog" mapstructure:"dog" lazy:"dog;foreign:id->dogs.owner_id"`
}

// Dog has many toys
type Dog struct {
	ID        uint      `gorm:"primarykey" json:"id" lazy:"id" mapstructure:"id"`
	CreatedAt time.Time `json:"created_at" lazy:"created_at" mapstructure:"created_at"`
	UpdatedAt time.Time `json:"updated_at" lazy:"updated_at" mapstructure:"updated_at"`
	Name      string    `json:"name" lazy:"name" mapstructure:"name"`
	OwnerID   uint      `json:"owner_id" lazy:"owner_id" mapstructure:"owner_id"`
	Toys      []Toy     `json:"toys" lazy:"toys" mapstructure:"toys"`
	Foods     []Food    `json:"foods" lazy:"foods" mapstructure:"foods" gorm:"many2many:dog_foods"`
	BreedID   uint      `json:"bread_id" lazy:"breed_id" mapstructure:"breed_id"`
}

// Profile belongs to a dog
type Profile struct {
	ID        uint      `gorm:"primarykey" json:"id" lazy:"id" mapstructure:"id"`
	CreatedAt time.Time `json:"created_at" lazy:"created_at" mapstructure:"created_at"`
	UpdatedAt time.Time `json:"updated_at" lazy:"updated_at" mapstructure:"updated_at"`
	Age       uint      `json:"age" lazy:"age" mapstructure:"age"`
	DogID     uint      `json:"-" lazy:"dog_id" mapstructure:"dog_id"`
	Dog       Dog       `json:"dog" lazy:"dog" mapstructure:"dog"`
}

type Toy struct {
	ID        uint      `gorm:"primarykey" json:"id" lazy:"id" mapstructure:"id"`
	CreatedAt time.Time `json:"created_at" lazy:"created_at" mapstructure:"created_at"`
	UpdatedAt time.Time `json:"updated_at" lazy:"updated_at" mapstructure:"updated_at"`
	Name      string    `json:"name" lazy:"name" mapstructure:"name"`
	DogID     uint      `json:"dog_id" lazy:"dog_id" mapstructure:"dog_id"`
}

type Breed struct {
	ID        uint      `gorm:"primarykey" json:"id" lazy:"id" mapstructure:"id"`
	CreatedAt time.Time `json:"created_at" lazy:"created_at" mapstructure:"created_at"`
	UpdatedAt time.Time `json:"updated_at" lazy:"updated_at" mapstructure:"updated_at"`
	Name      string    `json:"name" lazy:"name" mapstructure:"name"`
}

var gormDB *gorm.DB

func TestMain(m *testing.M) {
	logrus.SetLevel(logrus.DebugLevel)
	setup()
	defer os.Exit(m.Run())
	defer func() {
		teardown()
	}()
}

func initTeseDB() {
	var err error
	os.Remove("./test.db")
	gormDB, err = gorm.Open("sqlite3", "./test.db")
	if err != nil {
		panic(err)
	}

	gormDB.AutoMigrate(&Dog{}, &Profile{}, &Breed{}, &Toy{}, &Food{}, &Owner{})

	register(&Dog{}, &Profile{}, &Breed{}, &Toy{}, &Food{}, &Owner{})

	ownerNames := []string{
		`Bowie`,
		`Lennon`,
		`Cobain`,
		`Cash`,
		`Hendrix`,
		`Reed`,
		`Mercury`,
		`Smith`,
		`Pop`,
	}

	breedNames := []string{
		`Golden Retriever`,
		`Husky`,
		`Labrador`,
		`Alaskan Malamute`,
		`Corgi`,
		`Border Collie`,
	}

	for _, v := range ownerNames {
		gormDB.Create(&Owner{Name: v})
	}
	for _, v := range breedNames {
		gormDB.Create(&Breed{Name: v})
	}

	pedigree := &Food{Brand: `Pedigree`}
	gormDB.Create(pedigree)
	purina := &Food{Brand: `Purina`}
	gormDB.Create(purina)
	diamond := &Food{Brand: `Diamond`}
	gormDB.Create(diamond)

	dog := &Dog{Name: "Charlie", OwnerID: 1, BreedID: 1}
	gormDB.Create(dog)
	gormDB.Model(dog).Association("Foods").Append(pedigree, purina, diamond)
	gormDB.Model(dog).Association("Toys").Append(Toy{Name: `Tug`, DogID: dog.ID}, Toy{Name: `Toss`, DogID: dog.ID})

	dog = &Dog{Name: "Max", OwnerID: 2, BreedID: 2}
	gormDB.Create(dog)
	gormDB.Model(dog).Association("Foods").Append(pedigree, purina, diamond)
	gormDB.Model(dog).Association("Toys").Append(Toy{Name: `Tug`, DogID: dog.ID}, Toy{Name: `Toss`, DogID: dog.ID})

	dog = &Dog{Name: "Buddy", OwnerID: 3, BreedID: 3}
	gormDB.Create(dog)
	gormDB.Model(dog).Association("Foods").Append(pedigree, purina, diamond)

	dog = &Dog{Name: "Oscar", OwnerID: 4, BreedID: 4}
	gormDB.Create(dog)
	gormDB.Model(dog).Association("Foods").Append(pedigree, purina, diamond)

	dog = &Dog{Name: "Milo", OwnerID: 5, BreedID: 5}
	gormDB.Create(dog)
	gormDB.Model(dog).Association("Foods").Append(pedigree, purina, diamond)

	dog = &Dog{Name: "Archie", OwnerID: 6, BreedID: 6}
	gormDB.Create(dog)
	gormDB.Model(dog).Association("Foods").Append(pedigree, purina, diamond)

	dog = &Dog{Name: "Ollie", OwnerID: 7, BreedID: 1}
	gormDB.Create(dog)
	gormDB.Model(dog).Association("Foods").Append(pedigree, purina, diamond)

	dog = &Dog{Name: "Toby", OwnerID: 8, BreedID: 2}
	gormDB.Create(dog)
	gormDB.Model(dog).Association("Foods").Append(pedigree, purina, diamond)

	dog = &Dog{Name: "Jack", OwnerID: 9, BreedID: 3}
	gormDB.Create(dog)
	gormDB.Model(dog).Association("Foods").Append(pedigree, purina, diamond)

	for k := range ownerNames {
		gormDB.Create(&Profile{Age: uint(1 + k%3), DogID: uint(k)})
	}
}

func setup() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{DisableColors: true})
}

func teardown() {
	gormDB.Close()
	os.Remove("./test.db")
}
