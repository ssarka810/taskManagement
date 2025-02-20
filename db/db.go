package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
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
	db *gorm.DB
	User
	Task
}
type Database interface{
	UserRegister(username, password string)error
	GetUserByUsername(username string)(*User,error)
	CreateTask(inputTask *Task)error
	UpdateTask(inputTask *Task)error
	DeleteTask(inputTask *Task)error
	GetTasks()(*[]Task,error)
	GetTaskByTaskId(taskId int32)(*Task,error)
}

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

func NewDatabase(db *gorm.DB)Database{
	return &DatabaseDetails{
		db: db,
	}
}

func (store *DatabaseDetails)UserRegister(username, password string)error{
	userdetails :=&User{
		Username: username,
		Password: password,
	}
	if err :=store.db.Create(userdetails).Error;err!=nil{
		return err
	}
	return nil
}
func (store *DatabaseDetails)GetUserByUsername(username string)(*User,error){
	userdetails :=&User{}
	if err :=store.db.Where("username =?",username).First(userdetails).Error;err!=nil{
		return userdetails,err
	}
	return userdetails,nil
}

func (store *DatabaseDetails)CreateTask(inputTask *Task)error{
	if err :=store.db.Create(inputTask).Error;err!=nil{
		return err
	}
	return nil
}

func (store *DatabaseDetails)UpdateTask(inputTask *Task)error{
	if err :=store.db.Save(inputTask).Error;err!=nil{
		return err
	}
	return nil
}
func (store *DatabaseDetails)DeleteTask(inputTask *Task)error{
	if err :=store.db.Delete(inputTask).Error;err!=nil{
		return err
	}
	return nil
}
func (store *DatabaseDetails)GetTasks()(*[]Task,error){
	tasks :=&[]Task{}
	if err :=store.db.Find(tasks).Error;err!=nil{
		return tasks,err
	}
	return tasks,nil
}
func (store *DatabaseDetails)GetTaskByTaskId(taskId int32)(*Task,error){
	tasks :=&Task{}
	if err :=store.db.Find(tasks).Error;err!=nil{
		return tasks,err
	}
	return tasks,nil
}

