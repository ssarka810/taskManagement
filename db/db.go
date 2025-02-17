package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)
var gormDb *gorm.DB

type User struct {
	ID       int32  `gorm:"primaryKey;autoIncrement"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}
type Task struct {
	ID          int32     `gorm:"primaryKey;autoIncrement"`
	UserID      int32     `gorm:"not null"`
	Title       string    `gorm:"not null"`
	Description string
	DueDate     time.Time
	Status      string    `gorm:"not null"`
}

type DatabaseDetails struct{

}
type Database interface{}

func Init(host,port string)(*gorm.DB,error){
	dns :=url.URL{
		User: url.UserPassword("postgres","postgres"),
		Scheme: "postgres",
		Host: fmt.Sprintf("%s:%s",host,port),
		Path: "postgres", //db, this is the default db
		RawQuery: (&url.Values{"sslmode":[]string{"disable"}}).Encode(),
	}
	//connect to the default db
	db, err :=sql.Open("postgres",dns.String())
	if err!=nil{
		return nil, fmt.Errorf("error connecting to postgres: %v", err)
	}
	defer db.Close()

	// checking if the task_management database exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", "task_management").Scan(&exists)
	if err != nil {
			return nil, fmt.Errorf("error checking database existence: %v", err)
	}

	if !exists{
		logrus.Info("creating database [ task_management ]")
		_, err:=db.Exec("CREATE DATABASE task_management")
		if err!=nil{
			return nil,err
		}
		logrus.Info("Database [ task_management ] created successfully.")
	}
	dns.Path="task_management"

	gormDb, err =gorm.Open("postgres",dns.String())
	if err!=nil{
		return nil, fmt.Errorf("error connecting to task_management: %v", err)
	}
	gormDb.SingularTable(true)
	gormDb.Debug().AutoMigrate(
		&User{},
		&Task{},
	)

	return gormDb,nil

}

func GetDB()*gorm.DB{
	return gormDb
}