package lazy

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Response struct {
	Data     interface{} `json:"data"`
	ErrorMsg string      `json:"error_msg"`
	ErrorNo  int         `json:"error_no"`
}

// Food many to many dog
type Food struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at" mapstructure:"created_at"`
	UpdatedAt time.Time `json:"updated_at" mapstructure:"updated_at"`
	Brand     string    `json:"brand"`
}

// Owner has one dog
type Owner struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at" mapstructure:"created_at"`
	UpdatedAt time.Time `json:"updated_at" mapstructure:"updated_at"`
	Name      string    `json:"name"`
	Dog       Dog       `json:"dog"`
}

// Dog has many toys
type Dog struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at" mapstructure:"created_at"`
	UpdatedAt time.Time `json:"updated_at" mapstructure:"updated_at"`
	Name      string    `json:"name"`
	OwnerID   uint      `json:"owner_id"`
	Toys      []Toy     `json:"toys"`
	Foods     []Food    `json:"foods" gorm:"many2many:dog_foods"`
	BreedID   uint      `json:"bread_id"`
}

// Profile belongs to a dog
type Profile struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at" mapstructure:"created_at"`
	UpdatedAt time.Time `json:"updated_at" mapstructure:"updated_at"`
	Age       uint      `json:"age"`
	DogID     uint      `json:"-"`
	Dog       Dog       `json:"dog"`
}

type Toy struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at" mapstructure:"created_at"`
	UpdatedAt time.Time `json:"updated_at" mapstructure:"updated_at"`
	Name      string    `json:"name"`
	DogID     uint      `json:"dog_id"`
}

type Breed struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at" mapstructure:"created_at"`
	UpdatedAt time.Time `json:"updated_at" mapstructure:"updated_at"`
	Name      string    `json:"name"`
}

var gormDB *gorm.DB

func TestMain(m *testing.M) {
	logrus.SetFormatter(&nested.Formatter{
		TrimMessages:    true,
		TimestampFormat: "0102-150405",
		NoFieldsSpace:   true,
		HideKeys:        false,
		ShowFullLevel:   true,
		CallerFirst:     false,
		FieldsOrder:     []string{"component", "category"},
		CustomCallerFormatter: func(f *runtime.Frame) string {
			s := strings.Split(f.Function, ".")
			funcName := s[len(s)-1]
			return fmt.Sprintf(" [%s:%d][%s()]", path.Base(f.File), f.Line, funcName)
		},
	})
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.DebugLevel)
	setup()
	defer os.Exit(m.Run())
	defer func() {
		teardown()
	}()
}

func initTestDB() {
	var err error
	os.Remove("./test.db")

	gormDB, err = gorm.Open(sqlite.Open("./test.db"), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				// LogLevel:      logger.Info,
				LogLevel: logger.Silent,
				Colorful: true,
			},
		),
	})

	if err != nil {
		panic(err)
	}

	gormDB.AutoMigrate(&Dog{}, &Profile{}, &Breed{}, &Toy{}, &Food{}, &Owner{})

	Register(&Dog{}, &Profile{}, &Breed{}, &Toy{}, &Food{}, &Owner{})

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
	// gormDB.Model(dog).Association("Toys").Append(Toy{Name: `Tug`, DogID: dog.ID}, Toy{Name: `Toss`, DogID: dog.ID})
	gormDB.Model(dog).Association("Toys").Append(&Toy{Name: `Tug`}, &Toy{Name: `Toss`})

	dog = &Dog{Name: "Max", OwnerID: 2, BreedID: 2}
	gormDB.Create(dog)
	gormDB.Model(dog).Association("Foods").Append(pedigree, purina, diamond)
	gormDB.Model(dog).Association("Toys").Append(&Toy{Name: `Tug`, DogID: dog.ID}, &Toy{Name: `Toss`, DogID: dog.ID})

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
	logrus.SetFormatter(&nested.Formatter{
		TrimMessages:    true,
		TimestampFormat: "-",
		NoFieldsSpace:   true,
		HideKeys:        false,
		ShowFullLevel:   true,
		CallerFirst:     false,
		FieldsOrder:     []string{"component", "category"},
		CustomCallerFormatter: func(f *runtime.Frame) string {
			s := strings.Split(f.Function, ".")
			funcName := s[len(s)-1]
			return fmt.Sprintf(" [%s:%d][%s()]", path.Base(f.File), f.Line, funcName)
		},
	})
	logrus.SetReportCaller(true)
}

func teardown() {
	d, _ := gormDB.DB()
	d.Close()
	os.Remove("./test.db")
}
