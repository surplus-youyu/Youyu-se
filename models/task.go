package models

import "time"

const (
	TaskTypeSurvey = "TASK_TYPE_SURVEY"

	TaskStateCreated  = "CREATED"
	TaskStatePending  = "PENDING"
	TaskStateFinished = "FINISHED"
)

type Task struct {
	ID          int     `gorm:"column:id" json:"id"`
	Title       string  `gorm:"column:title" json:"title"`
	Creator     int     `gorm:"column:creator" json:"creator"`
	Reward      float32 `gorm:"column:reward" json:"reward"`
	Type        string  `gorm:"column:type" json:"type"`
	Description string  `gorm:"column:description" json:"description"`
	Content     string  `gorm:"column:content" json:"content"`
	State       string  `gorm:"column:state" json:"state"`
	// AssignedTo  int       `gorm:"column:assigned_to" json:"assigned_to"` // if 0 is unassigned
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (t *Task) TableName() string {
	return "task"
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
	if err := DB.Find(&user, User{Uid: user.Uid}); err != nil {
		tx.Rollback()
		panic(err)
	}
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
