package tools

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"
)

// IsPathExist 判断文件是否存在
func IsPathExist(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		return false
	}

	return true
}

// IsEmptyDir 判断文件夹是否为空，传入文件时返回false，传入空路径返回true
func IsEmptyDir(dir string) bool {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return true
	}
	if info.IsDir() {
		files, _ := ioutil.ReadDir(dir)
		return len(files) == 0
	}
	return false
}

// IsDir 判断是否是文件夹
func IsDir(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// EnsurePathExist 确保路径文件夹存在
func EnsurePathExist(filename string) error {
	if !IsPathExist(filename) {
		return os.MkdirAll(filename, 0777)
	}
	return nil
}

// EmptyToString 输入为空字符串时，用给定字符串替换
func EmptyToString(src, val string) string {
	if src == "" {
		return val
	}
	return src
}

// MD5 计算字符串md5
func MD5(elem ...string) string {
	md5h := md5.New()
	for _, el := range elem {
		if el != "" {
			md5h.Write([]byte(el))
		}
	}
	return hex.EncodeToString(md5h.Sum([]byte("guomio")))
}

//GetUUID get uuid
func GetUUID() string {
	u := uuid.NewV1()
	return u.String()
}

// Join filepath.Join
func Join(elem ...string) string {
	return filepath.Join(elem...)
}

// CombineURLs 组合url
func CombineURLs(elem ...string) string {
	if len(elem) == 0 {
		return ""
	}
	baseURL := elem[0]

	urls := []string{strings.Trim(baseURL, "/")}
	for _, u := range elem[1:] {
		urls = append(urls, strings.Trim(u, "/"))
	}
	url := strings.Join(urls, "/")
	url = strings.ReplaceAll(url, "/?", "?")
	url = strings.ReplaceAll(url, "=/", "=")
	url = strings.ReplaceAll(url, "/&", "&")
	return url
}

// IntToString int转string
func IntToString(n int) string {
	return strconv.Itoa(n)
}

// StringToInt string转int，转换错误时返回 0
func StringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
