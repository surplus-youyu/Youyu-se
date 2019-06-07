package models

type Task struct {
	Tid       int    `gorm:"column:tid;AUTO_INCREMENT"`
	Owner     int    `gorm:"column:owner"`
	Title     string `gorm:"column:title"`
	Summary   string `gorm:"column:content"`
	Type      string `gorm:"column:type"`
	Bounty    int    `gorm:"column:bounty"`
	Extra     string `gorm:"column:extra"`
	Enclosure string `gorm:"column:enclosure"`
}

func GetTaskById(tid int32) []Task {
	var result []Task
	DB.Table("task").Where("id=?", tid).Find(&result)
	return result
}

func CreateNewTask(task Task) {
	DB.Table("task").Create(&task)
}

func GetAllTasks() []Task {
	var result []Task
	DB.Table("task").Find(&result)
	return result
}

func GetSurveyByPublisherId(pid int32) []Task {
	var result []Task
	DB.Table("task").Where("owner=?", pid).Find(&result)
	return result
}
