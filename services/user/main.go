package main

import (
	_ "database/sql"
	"encoding/json"
	"errors"
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

		//TODO: don't expose the database structure (use a Rep)
		body := []UserRow{
			UserRow{
				GivenName:  "John",
				FamilyName: "Doe",
			},
			UserRow{
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

// Row contains common columns for all tables.
type Row struct {
	ID      []byte      `json:"_"`
	UUID    satori.UUID `json:"id" gorm:"-"`
	Created time.Time   `json:"created"`
	Updated time.Time   `json:"updated"`
	Deleted *time.Time  `json:"deleted"`
}

// UserRow is the model for the user table
type UserRow struct {
	Row
	EmailAddress string `json:"emailAddress"`
	GivenName    string `json:"givenName"`
	FamilyName   string `json:"familyName"`
	Status       Status `json:"status"`
	Role         Role   `json:"role"`
	password     []byte
	PictureURL   string `json:"pictureUrl"`
	Version      int    `json:"version"`
}

// TableName sets the table name for UserRow
func (UserRow) TableName() string {
	return "users"
}

type wrapper struct {
	Value *string
}

// BeforeCreate will populate the timestamps
func (row *Row) BeforeCreate(scope *gorm.Scope) error {

	var wrapped = wrapper{}

	scope.DB().Raw("SELECT BIN_TO_UUID(UUID_TO_BIN(UUID(),true)) as value FROM dual").Scan(&wrapped)

	id, err := satori.FromString(*wrapped.Value)
	if err != nil {
		log.Printf("BeforeCreate() err=%v", err)
	}

	row.ID = id.Bytes()
	row.Created = time.Now()
	row.Updated = time.Now()
	return nil
}

// AfterCreate wil populate the UUID
func (row *Row) AfterCreate(scope *gorm.Scope) error {

	uuid, err := satori.FromBytes(row.ID)
	if err != nil {
		log.Printf("AfterCreate() err=%v", err)
	}

	row.UUID = uuid

	return nil
}

// BeforeUpdate will populate the timestamps
func (row *Row) BeforeUpdate(scope *gorm.Scope) error {

	row.Updated = time.Now()
	return nil
}

// AfterUpdate checks the rows affected
func (row *Row) AfterUpdate(scope *gorm.Scope) error {

	log.Printf("AfterUpdate() db.RowsAffected=%d", scope.DB().RowsAffected)

	if scope.DB().RowsAffected < 1 {
		return errors.New("AfterUpdate.NoRowsAffected")
	}

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

	db.LogMode(true)

	user := &UserRow{
		EmailAddress: "john@doe.com",
		Status:       StatusInvited,
		Role:         RoleUser,
	}

	// Create
	db.Create(user)

	// optimistic locking
	db.Model(user).Where("version = ?", user.Version).Updates(UserRow{EmailAddress: "2222@doe.com", Version: user.Version + 1})

	log.Printf("other() user=%#v", user)

}
