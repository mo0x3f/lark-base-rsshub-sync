package utils

import (
	"os"
)

// WriteFile 将字节集数据写入到指定文件中，并返回错误信息
//
// file 字符串参数，传入写入文件路径，
// data 字节集参数，传入写入的数据。
func WriteFile(file string, data []byte) error {
	// 写文件
	return os.WriteFile(file, data, 0644)
}

// ReadFile 读取文件
//
// file 字符串参数，传入文件路径
func ReadFile(file string) ([]byte, error) {
	b, err := os.ReadFile(file)

	if err != nil {
		return nil, err
	}

	return b, nil
}

// Exists 检查文件是否存在
func Exists(file string) bool {
	// 获取文件信息
	_, err := os.Stat(file)
	// 检查错误
	if err == nil {
		return true
	}
	// 是否不存在
	if os.IsNotExist(err) {
		return false
	}

	return false
}
