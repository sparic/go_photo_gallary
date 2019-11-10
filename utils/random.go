package utils

import (
	"bytes"
	"math/rand"
	"strconv"
	"time"
)

func GetRandomNum(num int) string {
	rand.Seed(time.Now().UnixNano())
	randStr := ""
	for i := 0; i < num; i++ {
		var buffer bytes.Buffer
		buffer.WriteString(randStr)
		buffer.WriteString(strconv.Itoa(rand.Intn(100)))
		randStr = buffer.String()
	}
	return randStr
}
