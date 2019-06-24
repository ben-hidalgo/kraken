package main

import (
	_ "database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

func jsonResponse(w http.ResponseWriter, body interface{}, status int) {

	j, err := json.Marshal(body)

	if err != nil {
		log.Printf("JsonResponse() err=%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(j)
}

type user struct {
	GivenName  string
	FamilyName string
}

func main() {

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))

		body := []user{
			user{
				GivenName:  "John",
				FamilyName: "Doe",
			},
			user{
				GivenName:  "Jane",
				FamilyName: "Doe",
			},
		}

		jsonResponse(w, body, 200)

	})

	other()

	// log.Printf("GET /users listening on 9092")
	// log.Fatal(http.ListenAndServe(":9092", nil))
}

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

// Base contains common columns for all tables.
type Base struct {
	ID        uuid.UUID  `gorm:"type:binary(16);primary_key;"`
	Version   int        `json:"version"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"update_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.NewV4()
	return scope.SetColumn("ID", uuid)
}

// User is the model for the user table.
type User struct {
	Base
	SomeFlag bool `gorm:"column:some_flag;not null;default:true" json:"some_flag"`
	// Profile s `json:"profile"`
}

func other() {
	db, err := gorm.Open("mysql", "local:local@/local?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Printf("other() open err=%s", err)
		return
	}
	defer db.Close()

	driver, err := mysql.WithInstance(db.DB(), &mysql.Config{})
	if err != nil {
		log.Printf("other() withInstance err=%s", err)
		return
	}

	m, err := migrate.NewWithDatabaseInstance("file://./migrations", "mysql", driver)
	if err != nil {
		log.Printf("other() migrate err=%#v", err)
		return
	}

	m.Steps(1)

	// db.LogMode(true)

	// Migrate the schema
	db.AutoMigrate(&User{})

	db.AutoMigrate(&Product{})

	// Create
	db.Create(&Product{Code: "L1212", Price: 1000})

	db.Create(&User{SomeFlag: true})
	return

	// Read
	var product Product
	db.First(&product, 1)                   // find product with id 1
	db.First(&product, "code = ?", "L1212") // find product with code l1212

	// Update - update product's price to 2000
	db.Model(&product).Update("Price", 2000)

	// Delete - delete product
	// db.Delete(&product)
}
