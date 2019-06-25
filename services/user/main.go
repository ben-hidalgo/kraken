package main

import (
	_ "database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
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

	createOne(db)

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {

		var users []UserRow

		db.Find(&users)

		jsonResponse(w, users, 200)

	})

	log.Printf("GET /users listening on 9092")
	log.Fatal(http.ListenAndServe(":9092", nil))
}

func createOne(db *gorm.DB) {

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

// Row contains common columns for all tables.
type Row struct {
	ID      []byte      `json:"-"`
	UUID    satori.UUID `json:"id" gorm:"-"`
	Created EpochTime   `json:"created"`
	Updated EpochTime   `json:"updated"`
	Deleted *EpochTime  `json:"deleted,omitempty"`
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

//EpochTime supports JSON formatting
type EpochTime time.Time

//MarshalJSON marshals EpochTime
func (et EpochTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(et).Unix(), 10)), nil
}

//UnmarshalJSON unmarshalls EpochTime
func (et *EpochTime) UnmarshalJSON(data []byte) error {
	t := strings.Trim(string(data), `"`) // Remove quote marks from around the JSON string
	sec, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		return err
	}
	epochTime := time.Unix(sec, 0)
	*et = EpochTime(epochTime)
	return nil
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
	row.Created = EpochTime(time.Now())
	row.Updated = EpochTime(time.Now())
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

	row.Updated = EpochTime(time.Now())
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
	StatusInvited    Status = "INVITED"
	StatusUnverified Status = "UNVERIFIED"
	StatusVerified   Status = "VERIFIED"
	StatusDisabled   Status = "DISABLED"
	StatusBanned     Status = "BANNED"
	StatusDeleted    Status = "DELETED"
)

// Role of a User
type Role string

// The User.Role constants
const (
	RoleAdmin    Role = "ADMIN"
	RoleCSM      Role = "CSM"
	RoleInternal Role = "INTERNAL"
	RoleCustomer Role = "CUSTOMER"
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
