package utils

import (
	"fmt"
	"log"
	"strconv"
)

func StrToInt64(str string) int64 {
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.Println("Invalid id")
	}

	return num
}

func Int64ToStr(num int64) string {
	return strconv.FormatInt(num, 10)
}

func StrIsEmpty(value string) bool {
	if value == "" {
		return false
	} else {
		return true
	}
}

func StrToInt32(str string) int32 {
	i, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
	}
	return int32(i)
}


