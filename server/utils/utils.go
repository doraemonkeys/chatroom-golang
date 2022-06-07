package utils

import (
	"fmt"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// StrEncrypt 对传入字符串进行加密
func StrEncrypt(str string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.MinCost)
	if err != nil {
		fmt.Println("in func StrEncrypt,GenerateFromPassword failed,err:", err)
		return ""
	}
	return string(hash)
}

// StrMatch 对传入的加密字符串进行比对,str2为明文
func StrMatch(str1 string, str2 string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(str1), []byte(str2))
	return err == nil
}

//对传入的用户ID进行基本的验证(字符|长度等)
func UserIDTest(userID string) bool {
	if userID == "" {
		return false
	}
	for _, v := range userID {
		if v >= '0' && v <= '9' {
			continue
		}
		return false
	}
	return true
}


// Str2uint64 将字符串转换为uint64 使用前请确保传入的字符串是合法的
func Str2uint64(str string) uint64 {
	res, _ := strconv.ParseUint(str, 10, 64)
	return res
}

// Str2int64 将字符串转换为int64 使用前请确保传入的字符串是合法的
func Str2int64(str string) int64 {
	res, _ := strconv.ParseInt(str, 10, 64)
	return res
}

// Str2int32 将字符串转换为int32 使用前请确保传入的字符串是合法的
func Str2int32(str string) int32 {
	res, _ := strconv.ParseInt(str, 10, 32)
	return int32(res)
}
// uint64ToStr 将uint64转换为字符串 使用前请确保传入的uint64是合法的
func Uint64ToStr(num uint64) string {
	return fmt.Sprintf("%v",num)
}


//对传入的密码进行基本的验证(字符|长度等)
func PasswordTest(password string) bool {
	if password == "" {
		return false
	}
	for i := 0; i < len(password); i++ {
		if password[i] == ' ' || password[i] == '\n' || password[i] == '\t' {
			return false
		}
	}
	return true
}

//对传入的用户名进行基本的验证(字符|长度等)
func UserNameTest(userName string) bool {
	if userName == "" {
		return false
	}
	for i := 0; i < len(userName); i++ {
		if userName[i] == ' ' || userName[i] == '\n' || userName[i] == '\t' {
			return false
		}
	}
	return true
}
