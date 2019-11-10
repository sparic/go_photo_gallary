package utils

import (
	"fmt"
	"strconv"
)

func ConvertString2Int(value string) int {
	intVal, err := strconv.Atoi(value)
	if err != nil {
		fmt.Println("err", err)
	}
	return intVal
}