/*
 * @Author       : lvyitao lvyitao@fanhaninfo.com
 * @Date         : 2024-06-07 15:43:28
 * @LastEditTime : 2024-06-07 16:13:03
 */
package cmd

import (
	"log"
	"os"
	"path/filepath"
	"src/utils"
	"strconv"
)

var basePath = utils.GetExcutePath()
var pid = -1
var pidFile string = filepath.Join(basePath,"pid")


func initDaemon(){
	ok,_ := utils.PathExists(pidFile)
	if ok {
		bytes, err := os.ReadFile(pidFile)
		if err != nil {
			log.Fatal("failed to read pid file", err)
		}
		id, err := strconv.Atoi(string(bytes))
		if err != nil {
			log.Fatal("failed to parse pid data", err)
		}
		pid = id
	}
}