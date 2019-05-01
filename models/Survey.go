package models

type Survey struct {
	Id          int    `gorm:"column:id"`
	PublisherId int    `gorm:"column:publisher_id"`
	Title       string `gorm:"column:title"`
	Content     string `gorm:"column:content"`
}

func GetSurveyById(sid int) []Survey{
	var result []Survey
	DB.Table("survey").Where("id=?", sid).Find(&result)
	return result
}