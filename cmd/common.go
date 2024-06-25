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
	"src/factory"
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
			log.Println("failed to read pid file", err)
			pid = -1
			return
		}
		id, err := strconv.Atoi(string(bytes))
		if err != nil {
			log.Println("failed to parse pid data", err)
			pid = -1
			return
		}
		if utils.IsProcessAlive(id){
			name,err := utils.GetProcessName(id)
			if err != nil{
				log.Println(err)
				pid = -1
				return
			}
			if factory.ProductName == name{
				pid = id
			}
		}else{
			pid = -1
		}
		
	}
}