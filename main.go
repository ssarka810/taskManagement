package main

import (
	"github.com/sirupsen/logrus"
	"github.com/ssarkar/taskMamagement/db"
	router "github.com/ssarkar/taskMamagement/handlers"
)


func main(){
	newDb , err :=db.Init("localhost","5432")
	if err!=nil{
		logrus.Panic("Unable to create DB connection. Error: ",err)
	}
	database :=db.NewDatabase(newDb)
	server :=router.NewServer(database)
	server.Init()
	server.Start("9091")
}