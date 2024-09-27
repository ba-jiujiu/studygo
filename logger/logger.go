package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Logger interface {
	Log(content string)
}

type Terminal struct {
}

type file struct {
	path string
}

func File(path string) *file {
	return &file{
		path: path,
	}
}

// Log 输出日志内容到控制台
func (t *Terminal) Log(content string) {
	cur := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%+v: %s", cur, content)
}

func (f *file) Log(content string) {
	// 转为绝对路径
	absPath, err := filepath.Abs(f.path)
	if err != nil {
		fmt.Println("无法打开此路径")
	}
	// 检查路径是否存在, 不存在则创建
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		dirPath := filepath.Dir(absPath)
		// 使用os.MkdirAll确保目录都存在
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			fmt.Println("创建目录失败")
			return
		}
	}
	newFile, err := os.OpenFile(absPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("无法打开文件")
		return
	}
	defer func(newFile *os.File) {
		err := newFile.Close()
		if err != nil {
			fmt.Println("文件关闭异常")
		}
	}(newFile)
	//追加文件内容
	content = time.Now().Format("2006-01-02 15:04:05") + ": " + content
	_, err = newFile.WriteString(content)
	if err != nil {
		fmt.Println("写出日志失败")
	}
}
