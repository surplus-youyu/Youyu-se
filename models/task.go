package models

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/surplus-youyu/Youyu-se/utils"
)

const (
	TaskTypeSurvey = "TASK_TYPE_SURVEY"
	TaskTypeCustom = "TASK_TYPE_CUSTOM"

	TaskStatusCreated  = "TASK_STATUS_CREATED"
	TaskStatusPending  = "TASK_STATUS_PENDING"
	TaskStatusFinished = "TASK_STATUS_FINISHED"

	AssignmentStatusPending = "ASSIGNMENT_STATUS_PENDING"
	AssignmentStatusJudging = "ASSIGNMENT_STATUS_JUDGING"
	AssignmentStatusSuccess = "ASSIGNMENT_STATUS_SUCCESS"
	AssignmentStatusFailed  = "ASSIGNMENT_STATUS_FAILED"
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
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
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

func IssueTaskRewards(task Task) {
	var creator User
	var assigneeList []User
	task.Status = TaskStatusFinished
	tx := DB.Begin()
	if err := tx.Find(&creator, User{Uid: task.Creator}).Error; err != nil {
		tx.Rollback()
		panic(err)
	}
	if err := tx.
		Table("task").
		Select("user.*").
		Joins("LEFT JOIN assignment ON task.id = assignment.task_id LEFT JOIN user ON assignment.assignee = user.uid").
		Where("task.id = ? AND assignment.status = ?", task.ID, AssignmentStatusSuccess).
		Find(&assigneeList).
		Error; err != nil {
		tx.Rollback()
		panic(err)
	}
	if len(assigneeList) == 0 {
		creator.Balance += task.Reward
		if err := tx.Save(&creator); err != nil {
			tx.Rollback()
			panic(err)
		}
		if err := tx.Save(&task).Error; err != nil {
			tx.Rollback()
			panic(err)
		}
		return
	}
	intReward := int(task.Reward * 100.0)
	extra := intReward % len(assigneeList)
	creator.Balance += float32(extra) / 100.0
	reward := intReward / len(assigneeList)
	for i := range assigneeList {
		assigneeList[i].Balance += float32(reward) / 100.0
	}
	if err := tx.Save(&task).Error; err != nil {
		tx.Rollback()
		panic(err)
	}
	if err := tx.Save(&creator).Error; err != nil {
		tx.Rollback()
		panic(err)
	}
	for i := range assigneeList {
		if err := tx.Save(&assigneeList[i]).Error; err != nil {
			tx.Rollback()
			panic(err)
		}
	}
	tx.Commit()
}

func UpsertTask(task Task) {
	if err := DB.Save(&task).Error; err != nil {
		panic(err)
	}
}

type Assgn struct {
	ID     int `json:"id"`
	TaskID int `json:"task_id"`
	// Assignee int    `json:"assignee"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetAssignmentByID(id int) Assignment {
	var assgn Assignment
	if err := DB.Find(&assgn, Assignment{ID: id}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			panic(utils.Error{404, "Not found", err})
		} else {
			panic(err)
		}
	}
	return assgn
}

func GetAssignmentListByUid(id int) []Assgn {
	var assgnList []Assgn
	if err := DB.
		Table("assignment").
		Select("assignment.id, assignment.task_id, assignment.created_at, assignment.updated_at, assignment.status, task.title, task.type").
		Joins("LEFT JOIN task ON assignment.task_id = task.id").
		Where("assignment.assignee = ?", id).
		Find(&assgnList).Error; err != nil {
		panic(err)
	}
	return assgnList
}

func GetAssignmentByTaskIDAndUid(taskID, uid int) Assignment {
	var assgn Assignment
	if err := DB.Find(&assgn, Assignment{TaskID: taskID, Assignee: uid}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			panic(utils.Error{404, "用户未认领任务", err})
		} else {
			panic(err)
		}
	}
	return assgn
}

func GetAssignmentListByTaskID(id int) []Assignment {
	var assignList []Assignment
	if err := DB.Find(&assignList, Assignment{TaskID: id}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return assignList
		} else {
			panic(err)
		}
	}
	return assignList
}

func UpsertAssignment(assgn Assignment) {
	assgn.UpdatedAt = time.Now()
	if err := DB.Save(&assgn).Error; err != nil {
		if strings.Contains(err.Error(), "Error 1062") {
			panic(utils.Error{403, "用户已认领过该任务", err})
		} else {
			panic(err)
		}
	}
}
