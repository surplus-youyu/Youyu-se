package models

type User struct {
	Uid      int     `gorm:"column:uid"`
	Password string	 `gorm:"column:password"`
	RealName string  `gorm:"column:real_name"`
	NickName string  `gorm:"column:nick_name"`
	Age      int     `gorm:"column:age"`
	Gender   string  `gorm:"column:gender"`
	Balance  float32 `gorm:"column:balance"`
	Major    string  `gorm:"column:major"`
	Grade    int     `gorm:"column:grade"`
	Phone    string  `gorm:"column:phone"`
	Email    string  `gorm:"column:email"`
}

func GetUserByEmail(email string) []User {
	var result []User
	DB.Table("user").Where("email=?", email).Find(&result)
	return result
}

func CreateNewUser(newUser User) {
	DB.Table("user").Create(&newUser)
}