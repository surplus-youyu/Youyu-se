package models

import (
	"github.com/jinzhu/gorm"
	"github.com/surplus-youyu/Youyu-se/utils"
)

type User struct {
	// Uid      int     `gorm:"column:uid"`
	Uid      int     `gorm:"primary_key"`
	Password string  `gorm:"column:password"`
	NickName string  `gorm:"column:nickname"`
	Balance  float32 `gorm:"column:balance"`
	Email    string  `gorm:"column:email"`
	Age      int     `gorm:"column:age"`
	Gender   string  `gorm:"column:gender"`
	Phone    string  `gorm:"column:phone"`
	Avatar   string  `gorm:"column:avatar"`
	Grade    string  `gorm:"column:grade"`
	Major    string  `gorm:"column:major"`
}

func (u User) TableName() string {
	return "user"
}

func GetUserByEmail(email string) []User {
	var result []User
	DB.Table("user").Where("email=?", email).Find(&result)
	return result
}

func GetUserById(id int) User {
	var user User
	if err := DB.Find(&user, User{Uid: id}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			panic(utils.Error{404, "用户不存在", err})
		} else {
			panic(err)
		}
	}
	return user
}

func CreateNewUser(newUser User) {
	DB.Table("user").Create(&newUser)
}

func UpdateUser(newUser User) {
	DB.Save(&newUser)
}
