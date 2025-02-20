package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/ssarkar/taskMamagement/db"
	"github.com/ssarkar/taskMamagement/middleware"
)

type Server struct{
	db db.Database
	router *mux.Router
}

func NewServer(db db.Database)*Server{
	return &Server{
		db: db,
	}
}

func (server *Server)Init(){
	server.setupRouter()
}

func (server *Server)setupRouter(){
	server.router=mux.NewRouter().StrictSlash(true)
	server.setupApi()
}

func (server *Server)setupApi(){
	router :=server.router
	router.Use(middleware.LoggermiddleWare)
	public :=router.PathPrefix("/v1").Subrouter()
	public.HandleFunc("/api/register",server.UserRegisterHandler).Methods("POST")
	public.HandleFunc("/api/login",server.LoginHandler).Methods("POST")

	// protected routers
	protected :=router.PathPrefix("/v1").Subrouter()
	protected.Use(middleware.AuthorizationMiddleware)
	protected.HandleFunc("/api/tasks",server.CreateTaskHandler).Methods("POST")
	protected.HandleFunc("/api/tasks/{task_id}",server.UpdateTaskHandler).Methods("PUT")
	protected.HandleFunc("/api/tasks/{task_id}",server.DeleteTaskHandler).Methods("DELETE")
	protected.HandleFunc("/api/tasks",server.GetTasksHandler).Methods("GET")
	protected.HandleFunc("/api/tasks/{task_id}/complete",server.MarkTaskAsCompleteHandler).Methods("PATCH")

}
func (server *Server)Start(port string){
	logrus.Info("starting server with port : ",port)
	logrus.Panic(http.ListenAndServe(fmt.Sprintf(":%s",port),server.router))
}