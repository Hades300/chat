package dao

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

type User struct {
	UserName string
	PassWord string
	gorm.Model
}

type UserDao struct{}

var db, err = gorm.Open("sqlite3", "user.db")

func init() {
	if err != nil {
		log.Fatal("Error connecting database ", err)
	}
}

func (u *UserDao) TableName() string {
	return "users"
}

func (u *UserDao) AddUser(user User) {
	db.Create(&user)
}

func (u *UserDao) FindUserById(user User) User {
	target := User{}
	db.Table("users").First(&target, user.ID)
	return target
}

func (u *UserDao) ValidateUserNameAndPassWord(user User) bool {
	count := 0
	db.Table("users").Model(&user).Count(&count)
	if count > 0 {
		return true
	} else {
		return false
	}
}

func (u *UserDao) CheckHaveUserName(userName string) bool {
	count := 0
	db.Table("users").Where("user_name = ?", userName).Count(&count)
	if count != 0 {
		return true
	} else {
		return false
	}
}
