/*
 * @Author       : lvyitao lvyitao@fanhaninfo.com
 * @Date         : 2024-05-30 09:03:25
 * @LastEditTime : 2024-06-07 16:35:17
 */
package utils

import (
	"context"
	"log"
	"net"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
)

func GetIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Println(err.Error())
	}
	defer conn.Close()
	ipAddress := conn.LocalAddr().(*net.UDPAddr)
	return ipAddress.IP.String()
}

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
