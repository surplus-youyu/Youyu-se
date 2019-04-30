package models

import (
	"fmt"
)

type User struct {
	Uid      string  `gorm:"column:uid"`
	Password string	 `gorm:"column:password"`
	RealName string  `gorm:"column:real_name"`
	NickName string  `gorm:"column:nick_name"`
	Age      string  `gorm:"column:age"`
	Gender   string  `gorm:"column:gender"`
	Balance  float32 `gorm:"column:balance"`
	Major    string  `gorm:"column:major"`
	Grade    int     `gorm:"column:grade"`
	Phone    string  `gorm:"column:phone"`
	Email    string  `gorm:"column:email"`
}

func GetUserByEmail(email string) User {
	result := User{}
	DB.Table("user").Where("email=?", email).Find(&result)
	fmt.Println(result)
	return result
}