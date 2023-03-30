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
	gorm.Model
	Username string `gorm:"primaryKey;not null;"`
	Password string 
}

type TodoElement struct{
	gorm.Model
	ID int `gorm:"primaryKey;autoincrement;"`
	Username  string
	User User `gorm:"foreignKey:Username;reference:Username;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Title string
	Done bool
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

	if err := Db.AutoMigrate(new(TodoElement)); err != nil{
		fmt.Println(err)
	}
	return Db
}

func CheckUserPass(db *gorm.DB, username, password string) bool{
	user := new(User)
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		log.Println(err)
		return false
	}
	log.Println(user.Username)
	log.Println(user.Password)
	return user.Password == password
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

func GetTodoList(db *gorm.DB, username string) []TodoElement{
	element := new([]TodoElement)
	db.Model(&TodoElement{}).Where("username = ?", username).Find(&element)
	return *element
}

func CreateTodo(db *gorm.DB, username, title string) error {
	element := new(TodoElement)
	element.Done = false
	element.Title = title
	element.Username = username
	return db.Create(element).Error
}

func DoneTodo(db *gorm.DB, id int, done bool) error {
	element := new(TodoElement)
	if err := db.Model(&TodoElement{}).Where("id=?", id).First(&element).Error; err != nil{
		return err
	}
	element.Done = done
	if err := db.Save(&element).Error; err != nil{
		return err
	}
	return nil
}

func DeleteTodo(db *gorm.DB, id int) error {
	return db.Where("id=?",id).Unscoped().Delete(&TodoElement{}).Error
}
