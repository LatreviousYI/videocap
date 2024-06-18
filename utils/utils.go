/*
 * @Author       : lvyitao lvyitao@fanhaninfo.com
 * @Date         : 2024-05-30 09:03:25
 * @LastEditTime : 2024-06-07 16:35:17
 */
package utils

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
)


func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

type Config struct{
	Port int `toml:"port"`
	Debug bool `toml:"debug"`
}

func PathExistsAndCreate(path string) error {
	ok,_:=PathExists(path)
	if ok {
		return nil
	}
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func GetExcutePath() string{
	if strings.Contains(os.Args[0],"go-build"){
		dir,err := os.Getwd()
		if err != nil{
			panic(err)
		}
		return dir
	}else{
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exePath := filepath.Dir(ex)
		return exePath
	}
}



func GetConfig() Config {
	basePath := GetExcutePath()
	var config Config
	toml.DecodeFile(filepath.Join(basePath,"config/config.toml"),&config)
	return config
}


func SafeRunGoRunGoroutineForever(wg *sync.WaitGroup, ctx context.Context, fn func(), interval int64) {
	wg.Add(1)
	defer wg.Done()
	defer ErrCatch()
	for {
		select {
		case <-ctx.Done():
			log.Println("结束运行")
			return
		default:
			fn()
			if interval != 0 {
				time.Sleep(time.Duration(interval) * time.Second)
			}
		}
	}
}

func ErrCatch() {
	if err := recover(); err != nil {
		log.Printf("panic: %v stack:%s \n", err, string(debug.Stack()))
	}
}

func SafeRunGoRunGoroutineOnce(wg *sync.WaitGroup, ctx context.Context, fn func()) {
	wg.Add(1)
	defer wg.Done()
	defer ErrCatch()
	fn()
}



func IsLegalString(str string) bool {
	if str == "" {
		return false
	}
	strList := strings.Split(str, "")
	strListLen := len(strList)
	match, err := regexp.MatchString("^[\u4e00-\u9fa50-9a-zA-Z]+$", strList[0])
	if err != nil {
		log.Println(err.Error())
		return false
	}
	if !match {
		return false
	}
	if strListLen > 2 {
		match, err = regexp.MatchString("^[0-9a-zA-Z_\u4e00-\u9fa5]+$", strings.Join(strList[1:],""))
		if err != nil {
			log.Println(err.Error())
			return false
		}
		if !match {
			return false
		}
	}
	return true
}

func Today() string{
	return time.Now().Format("2006-01-02")
}

func DateTime() string{
	return time.Now().Format("2006_01_02_15_04_05")
}

func UnixTime() string{
	return fmt.Sprintf("%d",time.Now().Unix())
}



func GetMd5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	re := h.Sum(nil)
	return fmt.Sprintf("%x", re)
}

