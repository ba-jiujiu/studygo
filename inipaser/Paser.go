package inipaser

import (
	"bufio"
	"os"
	"reflect"
	"strings"
)

type ServerConfig struct {
	Port string `ini:"port"`
	Host string `ini:"host"`
}

type DatabaseConfig struct {
	Port     string `ini:"port"`
	Host     string `ini:"host"`
	Username string `ini:"username"`
	Password string `ini:"password"`
}

type Config struct {
	ServerConfig   `ini:"server"`
	DatabaseConfig `ini:"database"`
}

func ParseINI(path string, config interface{}) error {
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	// 释放资源
	defer file.Close()
	// 通过反射获取传入结构体的指针
	v := reflect.ValueOf(config).Elem()

	// 定义当前处理的section
	var curSection reflect.Value

	scanner := bufio.NewScanner(file)
	// 每次读取一行文件
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 忽略空行及注释
		if line == "" || strings.HasSuffix(line, ";") || strings.HasSuffix(line, "#") {
			continue
		}

		// 检查是否是节
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			sectionName := strings.TrimSpace(line[1 : len(line)-1])

			// 根据节名获取对应的字段
			for i := 0; i < v.NumField(); i++ {
				// reflect.Type才能调Field方法
				filed := v.Type().Field(i)
				if filed.Tag.Get("ini") == sectionName {
					curSection = v.Field(i)
					break
				}
			}
			continue
		}
		// 如果当前Section存在
		if curSection.IsValid() && strings.Contains(line, "=") {
			pair := strings.SplitN(line, "=", 2)
			key, value := strings.TrimSpace(pair[0]), strings.TrimSpace(pair[1])

			// 利用反射设置字段
			for i := 0; i < curSection.NumField(); i++ {
				filed := curSection.Type().Field(i)
				if filed.Tag.Get("ini") == key {
					curSection.Field(i).SetString(value)
					break
				}
			}
		}
		if err := scanner.Err(); err != nil {
			return err
		}
	}
	return nil
}
