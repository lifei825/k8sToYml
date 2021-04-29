package myUtil

import (
	"fmt"
	"os"
	"time"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func DateTime() string {
	currentTime := time.Now()
	datetime := currentTime.Format("2006-01-02-15-04-05")
	return datetime
}

func FileCreate(filename string) (bool, error) {
	f, err := os.Stat(filename)
	if err == nil {
		if f.IsDir() {
			fmt.Printf("文件创建失败：此路径为目录，请指定文件\n")
			return false, err
		} else if os.IsExist(err) {
			println("文件存在")
			return true, err
		} else {
			return true, err
		}
	}
	if os.IsNotExist(err) {
		println("文件不存在, 开始创建...")
		fp, err := os.Create(filename)
		if err != nil {
			fmt.Printf("文件创建失败：%s\n", err)
		}
		defer func() {_=fp.Close()}()
		return true, err
	}
	fmt.Printf("文件创建失败：%s\n", err)
	return false, err
}