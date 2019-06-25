package main

import (
	_ "database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jinzhu/gorm"
	satori "github.com/satori/go.uuid"
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
	ID      []byte      `json:"_"`
	UUID    satori.UUID `json:"id" gorm:"-"`
	Version int         `json:"version"`
	Created time.Time   `json:"created"`
	Updated time.Time   `json:"updated"`
	Deleted *time.Time  `json:"deleted"`
}

type wrapper struct {
	Value *string
}

// BeforeCreate will populate the timestamps
func (base *Base) BeforeCreate(scope *gorm.Scope) error {

	var wrapped = wrapper{}

	scope.DB().Raw("SELECT BIN_TO_UUID(UUID_TO_BIN(UUID(),true)) as value FROM dual").Scan(&wrapped)

	id, err := satori.FromString(*wrapped.Value)
	if err != nil {
		log.Printf("BeforeCreate() err=%v", err)
	}

	base.ID = id.Bytes()

	base.Created = time.Now()
	base.Updated = time.Now()
	return nil
}

// AfterCreate wil populate the UUID
func (base *Base) AfterCreate(scope *gorm.Scope) error {

	uuid, err := satori.FromBytes(base.ID)
	if err != nil {
		log.Printf("AfterCreate() err=%v", err)
	}

	base.UUID = uuid

	return nil
}

// BeforeUpdate will populate the timestamps
func (base *Base) BeforeUpdate(scope *gorm.Scope) error {
	base.Updated = time.Now()
	return nil
}

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

// MyLogger implements migrate.Logger
type MyLogger struct {
}

// Printf implementation of migrate.Logger.Printf
func (ml *MyLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v)
}

// Verbose implementation of migrate.Logger.Verbose
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

	// db.LogMode(true)

	user := &User{
		EmailAddress: "john@doe.com",
		Status:       StatusInvited,
		Role:         RoleUser,
	}

	// Create
	db.Create(user)

	log.Printf("other() ... user=%#v", user)

}
