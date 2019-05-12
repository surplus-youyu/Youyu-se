package models

type Survey struct {
	Id          int    `gorm:"column:id;AUTO_INCREMENT"`
	PublisherId int    `gorm:"column:publisher_id"`
	Title       string `gorm:"column:title"`
	Content     string `gorm:"column:content"`
}

func GetSurveyById(sid int32) []Survey {
	var result []Survey
	DB.Table("survey").Where("id=?", sid).Find(&result)
	return result
}

func CreateNewSurvey(survey Survey) {
	DB.Table("survey").Create(&survey)
}

func GetAllSurvey() []Survey {
	var result []Survey
	DB.Table("survey").Find(&result)
	return result
}

func GetSurveyByPublisherId(pid int32) []Survey {
	var result []Survey
	DB.Table("survey").Where("publisher_id=?", pid).Find(&result)
	return result
}