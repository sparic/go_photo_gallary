package service

import (
	"crypto/md5"
	"fmt"
	"io"
	"scratch_maker_server/models"
	"strconv"
)

func GetUserList(pageNum, pageSize, nameQuery string) []models.User {
	models.Db.LogMode(true)
	var userList []models.User
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		fmt.Println("err", err)
	}

	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		fmt.Println("err", err)
	}

	models.Db.Where("user_name LIKE ?", "%"+nameQuery+"%").Order("id asc").
		Offset((pageNumInt - 1) * pageSizeInt).Limit(pageSize).Find(&userList)
	return userList
}

func GetUserCount(nameQuery string) int {
	var count int
	db := models.Db.Model(&models.User{})
	if nameQuery != "" {
		db = db.Where("user_name LIKE ?", "%"+nameQuery+"%")
	}
	db.Count(&count)
	return count
}

func UpdateUser(user models.User) {
	models.Db.LogMode(true)
	//更新
	hash := md5.New()
	io.WriteString(hash, user.Password) // for safety, don't just save the plain text
	user.Password = fmt.Sprintf("%x", hash.Sum(nil))

	models.Db.Model(&user).Updates(user)
}

func DeleteUser(id uint) {
	var user = models.User{}
	user.ID = id
	models.Db.LogMode(true)
	models.Db.Delete(&user)
}
