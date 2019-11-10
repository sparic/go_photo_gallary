package models

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type User struct {
	BaseModel
	UserName string         `json:"userName" gorm:"column:user_name;type:varchar(16)"`
	Password string         `json:"password" gorm:"column:password;type:varchar(16)"`
	Email    string         `json:"email" gorm:"column:email;type:varchar(128)"`
	NickName string         `json:"nickName" gorm:"column:nick_name;type:varchar(128)"`
	Sex      bool           `json:"sex" gorm:"column:sex;type:tinyint(1)"`
	Birthday mysql.NullTime `json:"birthday" gorm:"column:birthday;type:datetime"`
}

type UserReq struct {
	ID       uint   `json:"id"`
	UserName string `json:"userName"`
	Password string `json:"password"`
	Email    string `json:"email"`
	NickName string `json:"nickName"`
	Sex      bool   `json:"sex"`
	Birthday string `json:"birthday" time_format:"2006-01-02" time_utc:"1"`
}

type UserPageResp struct {
	PageNum  int    `json:"pageNum"`
	PageSize int    `json:"pageSize"`
	Count    int    `json:"count"`
	Rows     []User `json:"rows"`
}

//TODO 这种方式很不优雅， 字符串转时间戳格式
func UserReq2User(userReq UserReq) User {
	var user User
	user.ID = userReq.ID
	user.UserName = userReq.UserName
	user.Password = userReq.Password
	user.Email = userReq.Email
	user.NickName = userReq.NickName
	user.Sex = userReq.Sex
	if userReq.Birthday != "" {
		birthD, err := time.Parse("2006-01-02 00:00:00", userReq.Birthday)
		if err != nil {
			log.Fatalln(err)
		}
		user.Birthday.Time = birthD
		user.Birthday.Valid = true
	}
	return user
}

var AuthExistsError = errors.New("auth already exists")

// 数据库插入User记录
func InsertUser(userEntity User) error {
	trx := Db.Begin()

	user := User{}
	trx.Set("gorm:query_option", "FOR UPDATE").
		Where("user_name = ?", userEntity.UserName).
		First(&user)
	if user.ID > 0 {
		return AuthExistsError
	}

	hash := md5.New()
	io.WriteString(hash, userEntity.Password) // for safety, don't just save the plain text
	user = userEntity
	user.Password = fmt.Sprintf("%x", hash.Sum(nil))

	err := trx.Create(&user).Error
	if err != nil {
		return err
	}

	// fmt.Println("insert success" + Test(100))

	defer trx.Commit()
	return nil
}

func Test(i int) string {
	var arr [10]int
	arr[i] = 123 //err panic
	//错误拦截必须配合defer使用  通过匿名函数使用
	defer func() {
		//恢复程序的控制权
		err := recover()
		if err != nil {
			fmt.Println("err happened!", err)
		}
	}()
	return "!!"
}

func SelectUserByName(userName string, trx *gorm.DB) User {
	var userByName User
	//校验用户名是否重复
	trx.Find(&userByName, "user_name = ?", userName)
	return userByName
}

func SelectUserById(id uint, trx *gorm.DB) User {
	var userById User
	//校验用户名是否重复
	trx.First(&userById, id)
	return userById
}

// Check if the auth is valid.
func CheckUser(username, password string) bool {
	trx := Db.Begin()
	defer trx.Commit()

	hash := md5.New()
	io.WriteString(hash, password)
	password = fmt.Sprintf("%x", hash.Sum(nil)) //	for safety, don't just save the plain text
	user := User{}
	trx.Set("gorm:query_option", "FOR UPDATE").
		Where("user_name = ? AND password = ?", username, password).
		First(&user)
	if user.ID > 0 {
		return true
	}
	return false
}
