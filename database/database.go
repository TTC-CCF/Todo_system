package database

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)


const (
	USERNAME = "popuku"
	PASSWORD = "123"
	NETWORK = "tcp"
	SERVER = "127.0.0.1"
	PORT = 3306
	DATABASE = "todo_system"
)

type User struct{
	Username string `json:"username" gorm:"primary_key"`
	Password string `json:""`
}



func ConnectDB() *gorm.DB{
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",USERNAME,PASSWORD,NETWORK,SERVER,PORT,DATABASE)
	Db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	
	if err != nil{
		fmt.Println(err)
	}

	if err := Db.AutoMigrate(new(User)); err != nil{
		fmt.Println(err)
	}
	return Db
}

func CheckUserPass(db *gorm.DB, username, password string) bool{
	user := new(User)
	user.Username = username
	if err := db.First(&user).Error; err != nil {
		return false
	}
	
	if user.Password != password {
		return false
	}
	return true
}

func EmptyUserPass(username, password string) bool {
	return strings.Trim(username, " ") == "" || strings.Trim(password, " ") == ""
}

func CheckUserExist(db *gorm.DB, username string) bool{
	count := int64(0)
	err := db.Model(&User{}).Where("username = ?", username).Count(&count).Error
	if err != nil{
		fmt.Println(err)
	}
	log.Println(count)
	return count > 0
}

func CreateUser(db *gorm.DB, username, password string) error {
	user := new(User)
	user.Password = password
	user.Username = username
	return db.Create(user).Error
}
