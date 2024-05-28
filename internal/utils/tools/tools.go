package tools

import (
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func RemoveDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func CheckPagination(pageSize int, pageIndex int) (int, int, int) {
	if pageSize <= 0 {
		pageSize = -1
		pageIndex = 0
	}
	// max page size is 100
	if pageSize > 100 {
		pageSize = 100
	}
	var offset int
	offset = pageSize * pageIndex
	if pageSize <= 0 || pageIndex < 0 {
		offset = -1
	}

	return pageSize, pageIndex, offset
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func CheckPassword(providedPassword, storedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(providedPassword))
	return err == nil
}

func StringToUint(str string) uint {
	value, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		panic(err)
	}
	return uint(value)
}
