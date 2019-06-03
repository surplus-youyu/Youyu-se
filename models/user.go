package models

type User struct {
	Uid      int     `gorm:"column:uid"`
	Password string  `gorm:"column:password"`
	NickName string  `gorm:"column:nickname"`
	Balance  float32 `gorm:"column:balance"`
	Email    string  `gorm:"column:email"`
	Age      int     `gorm:"column:age"`
	Gender   string  `gorm:"column:gender"`
	Phone    string  `gorm:"column:phone"`
	Avatar   string  `gorm:"column:avatar"`
}

func GetUserByEmail(email string) []User {
	var result []User
	DB.Table("user").Where("email=?", email).Find(&result)
	return result
}

func CreateNewUser(newUser User) {
	DB.Table("user").Create(&newUser)
}

func UpdateUser(newUser User) {
	DB.Save(&newUser)
}
