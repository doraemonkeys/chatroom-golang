package utils

import (
	"io"
	"os"
	"strconv"
)

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

//从文件末尾按行读取文件。
//name:文件路径 lineNum:读取行数(超过文件行数则读取全文)
func ReverseRead(name string, lineNum uint) ([]string, error) {
	//打开文件
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	//获取文件大小
	fs, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fs.Size()

	var offset int64 = -1   //偏移量，初始化为-1，若为0则会读到EOF
	char := make([]byte, 1) //用于读取单个字节
	lineStr := ""           //存放一行的数据
	buff := make([]string, 0, 100)
	for (-offset) <= fileSize {
		//通过Seek函数从末尾移动游标然后每次读取一个字节，offset为偏移量
		file.Seek(offset, io.SeekEnd)
		_, err := file.Read(char)
		if err != nil {
			return buff, err
		}
		if char[0] == '\n' {
			//判断文件类型为unix(LF)还是windows(CRLF)
			file.Seek(-2, io.SeekCurrent) //io.SeekCurrent表示游标放置于当前位置，逆向偏移2个字节
			//读完一个字节后游标会自动正向偏移一个字节
			file.Read(char)
			if char[0] == '\r' {
				offset-- //windows跳过'\r'
			}
			lineNum-- //到此读取完一行
			buff = append(buff, lineStr)
			lineStr = ""
			if lineNum == 0 {
				return buff, nil
			}
		} else {
			lineStr = string(char) + lineStr
		}
		offset--
	}
	buff = append(buff, lineStr)
	return buff, nil
}
