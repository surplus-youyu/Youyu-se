package models

type Answer struct {
	Id          int    `gorm:"column:id;AUTO_INCREMENT"`
	IntervieweeId int    `gorm:"column:interviewee_id"`
	Content     string `gorm:"column:content"`
}

func CreateNewAnswer(answer Answer) {
	DB.Table("answer").Create(&answer)
}
