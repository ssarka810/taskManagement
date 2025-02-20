package models

import (
	"time"

	"github.com/ssarkar/taskMamagement/db"
)

type TaskDetails struct{
	Title string `json:"Title,omitempty" validate:"required" `
	Description string `json:"description" `
	DueDate time.Time `json:"due_date,omitempty" validate:"required" `
	TaskStatus string `json:"status"`
}

func ConvertDbTasksForResponse(inputTasks []db.Task)([]TaskDetails){
	tasks :=[]TaskDetails{}
	for _,dbtask :=range inputTasks{
		task :=TaskDetails{}
		task.Description=dbtask.Description
		task.Title=dbtask.Title
		task.TaskStatus=dbtask.Status
		task.DueDate=dbtask.DueDate
		tasks=append(tasks, task)
	}
	return tasks
}

func ConvertInputTaskToDbTask(task *TaskDetails,userId int32)db.Task{
	dbTask :=db.Task{}
	dbTask.UserID=userId
	dbTask.Description=task.Description
	dbTask.DueDate=task.DueDate
	dbTask.Title=task.Title
	return dbTask
}