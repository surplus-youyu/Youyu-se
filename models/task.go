package models

import (
	"strings"
	"time"

	"github.com/surplus-youyu/Youyu-se/utils"
)

const (
	TaskTypeSurvey = "TASK_TYPE_SURVEY"
	TaskTypeCustom = "TASK_TYPE_CUSTOM"

	TaskStatusCreated  = "TASK_STATUS_CREATED"
	TaskStatusPending  = "TASK_STATUS_PENDING"
	TaskStatusFinished = "TASK_STATUS_FINISHED"

	AssignmentStatusPending = "PENDING"
	AssignmentStatusJudging = "JUDGING"
	AssignmentStatusSuccess = "SUCCESS"
	AssignmentStatusFailed  = "FAILED"
)

type Task struct {
	ID          int     `gorm:"column:id" json:"id"`
	Title       string  `gorm:"column:title" json:"title"`
	Creator     int     `gorm:"column:creator" json:"creator"`
	Reward      float32 `gorm:"column:reward" json:"reward"`
	Type        string  `gorm:"column:type" json:"type"`
	Limit       int     `gorm:"column:limit" json:"limit"`
	Assigned    int     `gorm:"column:assigned" json:"assigned"`
	Description string  `gorm:"column:description" json:"description"`
	// Demand      string       `gorm:"column:demand" json:"demand"`
	Content     string       `gorm:"column:content" json:"content"`
	Status      string       `gorm:"column:status" json:"status"`
	CreatedAt   time.Time    `gorm:"column:created_at" json:"created_at"`
	Assignments []Assignment `gorm:"ForeignKey:TaskID" json:"-"`
}

type Assignment struct {
	ID        int       `gorm:"primary_key" json:"id"`
	TaskID    int       `gorm:"column:task_id" json:"task_id"`
	Assignee  int       `gorm:"column:assignee" json:"assignee"`
	Status    string    `gorm:"column:status" json:"status"`
	Payload   string    `gorm:"column:payload" json:"payload"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (t *Task) TableName() string {
	return "task"
}

func (a *Assignment) TableName() string {
	return "assignment"
}

func GetTaskList() []Task {
	var taskList []Task
	if err := DB.Find(&taskList).Error; err != nil {
		panic(err)
	}
	return taskList
}

func CreateTask(task Task, user User) int {
	tx := DB.Begin()
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> Update
	if err := DB.Find(&user, User{Uid: user.Uid}); err != nil {
=======
	if err := DB.Find(&user, User{Uid: user.Uid}).Error; err != nil {
>>>>>>> Add: task assignment API
		tx.Rollback()
		panic(err)
	}
<<<<<<< HEAD
	user.Balance -= task.Reward
	if user.Balance < 0 {
		tx.Rollback()
		panic("余额不足")
	}
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		panic(err)
	}
	if err := tx.Create(&task).Error; err != nil {
		tx.Rollback()
=======
=======
>>>>>>> Update
	user.Balance -= task.Reward
	if user.Balance < 0 {
		tx.Rollback()
		panic("余额不足")
	}
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		panic(err)
	}
	if err := tx.Create(&task).Error; err != nil {
<<<<<<< HEAD
>>>>>>> Add: task
=======
		tx.Rollback()
>>>>>>> Update
		panic(err)
	}
	tx.Commit()
	return task.ID
}

func GetTaskByID(id int) Task {
	var task Task
	if err := DB.Find(&task, Task{ID: id}).Error; err != nil {
		panic(err)
	}
	return task
}

func UpsertTask(task Task) {
	if err := DB.Save(&task).Error; err != nil {
		panic(err)
	}
}

type Assgn struct {
	TaskID   int    `json:"task_id"`
	Assignee int    `json:"assignee"`
	Status   string `json:"status"`
	Title    string `json:"title"`
}

func GetAssignmentListByUid(id int) []Assgn {
	var assgnList []Assgn
	if err := DB.
		Table("assignment").
		Select("assignment.task_id, assignment.assignee, assignment.status, task.title").
		Joins("LEFT JOIN task ON assignment.task_id = task.id").
		Where("assignment.assignee = ?", id).
		Find(&assgnList).Error; err != nil {
		panic(err)
	}
	return assgnList
}

func UpsertAssignment(assgn Assignment) {
	if err := DB.Save(&assgn).Error; err != nil {
		if strings.Contains(err.Error(), "Error 1062") {
			panic(utils.Error{403, "用户已认领过该任务", err})
		} else {
			panic(err)
		}
	}
}
