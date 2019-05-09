package models

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/surplus-youyu/Youyu-se/config"
)

var (
	DB *gorm.DB
)

func init() {
	var error error
	DB, error = getDB()
	if error != nil {
		panic(error)
	}
}

func getDB() (*gorm.DB, error) {
	config := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", config.DBUser,
		config.DBPwd, config.DBHost, config.DBName)
	fmt.Println(config)
	db, err := gorm.Open("mysql", config)
	return db, err
}
