package main

import (
	_ "database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	// _ "github.com/go-sql-driver/mysql"
	// _ "github.com/golang-migrate/migrate/source/file"

	// "github.com/golang-migrate/migrate/v4"
	// "github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jinzhu/gorm"
	// uuid "github.com/satori/go.uuid"
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

func main() {

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {

		body := []User{
			User{
				GivenName:  "John",
				FamilyName: "Doe",
			},
			User{
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

// type Product struct {
// 	gorm.Model
// 	Code  string
// 	Price uint
// }

// User is the model for the user table
type User struct {
	Base
	EmailAddress string `json:"emailAddress"`
	GivenName    string `json:"givenName"`
	FamilyName   string `json:"familyName"`
	Status       Status `json:"status"`
	Role         Role   `json:"role"`
	password     []byte
	PictureURL   string `json:"pictureUrl"`
}

// Base contains common columns for all tables.
type Base struct {
	//ID        uuid.UUID  `gorm:"type:binary(16);primary_key;"`
	Version   int        `json:"version"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"update_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

// BeforeCreate will set a UUID rather than numeric ID.
// func (base *Base) BeforeCreate(scope *gorm.Scope) error {
// 	uuid := uuid.NewV4()
// 	return scope.SetColumn("ID", uuid)
// }

// Status of a User
type Status string

// The User.Status constants
const (
	StatusInvited  Status = "INVITED"
	StatusAccepted Status = "ACCEPTED"
	StatusVerified Status = "VERIFIED"
	StatusDeleted  Status = "DELETED"
	StatusDisabled Status = "DISABLED"
	StatusBanned   Status = "BANNED"
)

// Role of a User
type Role string

// The User.Role constants
const (
	RoleAdmin    Role = "ADMIN"
	RoleCSM      Role = "CSM"
	RoleEmployee Role = "EMPLOYEE"
	RoleUser     Role = "USER"
)

type MyLogger struct {
}

func (ml *MyLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v)
}
func (ml *MyLogger) Verbose() bool {
	return true
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

	m.Log = &MyLogger{}

	m.Up()
	//TODO: this doesn't actually fail because it is asynch; check the schema_migrations.dirty == 1?
	if err != nil {
		log.Printf("other() up err=%#v", err)
		return
	}
	log.Println("other() returning")
	return

	// db.LogMode(true)

	// Migrate the schema
	// db.AutoMigrate(&User{})
	// db.AutoMigrate(&Product{})

	// Create
	// db.Create(&Product{Code: "L1212", Price: 1000})

	// db.Create(&User{SomeFlag: true})
	return

	// Read
	// var product Product
	// db.First(&product, 1)                   // find product with id 1
	// db.First(&product, "code = ?", "L1212") // find product with code l1212

	// // Update - update product's price to 2000
	// db.Model(&product).Update("Price", 2000)

	// Delete - delete product
	// db.Delete(&product)
}
