package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/ssarkar/taskMamagement/models"
)

func (server *Server)CreateTaskHandler(w http.ResponseWriter, r *http.Request){
	taskDetails :=&models.TaskDetails{}
	if err:=json.NewDecoder(r.Body).Decode(taskDetails);err!=nil{
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	if err:=ValidateInputStructs(taskDetails);err!=nil{
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	username, ok :=r.Context().Value("username").(string)
	if !ok{
		http.Error(w,"Uname to get the username from the context",http.StatusInternalServerError)
		return
	}
	logrus.Info("username ",username)
	user,err:=server.db.GetUserByUsername(username)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	dbTask :=models.ConvertInputTaskToDbTask(taskDetails,user.ID)
	dbTask.Status="not-started"
	if err:=server.db.CreateTask(&dbTask);err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	w.Write([]byte("task is added successfully"))
}

func (server *Server)UpdateTaskHandler(w http.ResponseWriter, r *http.Request){
	taskDetails :=&models.TaskDetails{}
	if err:=json.NewDecoder(r.Body).Decode(taskDetails);err!=nil{
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	if err:=ValidateInputStructs(taskDetails);err!=nil{
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	task_id_input :=mux.Vars(r)["task_id"]
	if task_id_input==""{
		http.Error(w,"task id can not be empty",http.StatusBadRequest)
		return
	}
	task_id, err :=strconv.Atoi(task_id_input)
	if err!=nil{
		http.Error(w,"provide the correct taskId",http.StatusBadRequest)
		return
	}
	dbTask, err :=server.db.GetTaskByTaskId(int32(task_id))
	if errors.Is(err,gorm.ErrRecordNotFound){
		http.Error(w,"there is no task with the task_id",http.StatusBadRequest)
		return
	}else if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}else{
		dbTask.Description=taskDetails.Description
		dbTask.Title=taskDetails.Title
		dbTask.DueDate=taskDetails.DueDate
		if err:=server.db.UpdateTask(dbTask);err!=nil{
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}
	}
	w.Write([]byte("task is updated successfully"))
}

func (server *Server)DeleteTaskHandler(w http.ResponseWriter, r *http.Request){
	taskDetails :=&models.TaskDetails{}
	if err:=json.NewDecoder(r.Body).Decode(taskDetails);err!=nil{
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	if err:=ValidateInputStructs(taskDetails);err!=nil{
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	task_id_input :=mux.Vars(r)["task_id"]
	if task_id_input==""{
		http.Error(w,"task id can not be empty",http.StatusBadRequest)
		return
	}
	task_id, err :=strconv.Atoi(task_id_input)
	if err!=nil{
		http.Error(w,"provide the correct taskId",http.StatusBadRequest)
		return
	}
	dbTask, err :=server.db.GetTaskByTaskId(int32(task_id))
	if errors.Is(err,gorm.ErrRecordNotFound){
		http.Error(w,"there is no task with the task_id",http.StatusBadRequest)
		return
	}else if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}else{
		if err:=server.db.DeleteTask(dbTask);err!=nil{
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}
	}
	w.Write([]byte("task is deleted successfully"))
}

func (server *Server)GetTasksHandler(w http.ResponseWriter, r *http.Request){
	dbTask, err :=server.db.GetTasks()
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	taskDetails :=models.ConvertDbTasksForResponse(*dbTask)
	response , err :=json.Marshal(taskDetails)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (server *Server)MarkTaskAsCompleteHandler(w http.ResponseWriter, r *http.Request){
	task_id_input :=mux.Vars(r)["task_id"]
	if task_id_input==""{
		http.Error(w,"task id can not be empty",http.StatusBadRequest)
		return
	}
	task_id, err :=strconv.Atoi(task_id_input)
	if err!=nil{
		http.Error(w,"provide the correct taskId",http.StatusBadRequest)
		return
	}
	dbTask, err :=server.db.GetTaskByTaskId(int32(task_id))
	if errors.Is(err,gorm.ErrRecordNotFound){
		http.Error(w,"there is no task with the task_id",http.StatusBadRequest)
		return
	}else if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}else{
		dbTask.Status="complete"
		if err:=server.db.UpdateTask(dbTask);err!=nil{
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}
	}
	w.Write([]byte("task is updated successfully as 'complete'"))
}
