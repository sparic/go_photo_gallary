package utils

import (
	"time"
)

var now = time.Now()

var year = now.Format("2006")
var month = now.Format("01")
var day = now.Format("02")
var hour = now.Format("15")
var min = now.Format("04")
var second = now.Format("05")

// 获取当前时间戳
func GetCurrentTimestamp() string {
	currentTimestamp := year + month + day + hour + min + second
	return currentTimestamp
}
