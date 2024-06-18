package utils

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"src/common"
	"strings"
	"syscall"
)


func GetDeviceList()([]string,error){
	var deviceList []string
	cmd := exec.Command("/bin/bash","-c",string("ls /dev/video*"))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout  // 标准输出
	cmd.Stderr = &stderr  // 标准错误
	err := cmd.Run()
	if err != nil {
		return deviceList,err
	}
	outStr, errStr := stdout.String(), stderr.String()
	if errStr != ""{
		return deviceList,errors.New(errStr)
	}
	deviceList = strings.Fields(outStr)
	return deviceList,nil
}

// ./mjpg_streamer -i "./input_uvc.so -n -f 30 -r 640x480 -d /dev/video0"  -o "./output_http.so -w ./www" &
func StartMjpgStreamer(resolutionRatio,deviceId string)error{
	basePath := GetExcutePath()
	cmd := exec.Command("/bin/bash","-c", fmt.Sprintf(`./mjpg_streamer -i "./input_uvc.so -n -f 30 -r %s -d %s"  -o "./output_http.so -w ./www"`,resolutionRatio,deviceId))
	cmd.Dir = filepath.Join(basePath,"mjpg-streamer/mjpg-streamer-experimental")
	err := cmd.Start()
	if err != nil {
		return err
	}
	common.EditMjpgPid(cmd.Process.Pid)
	return nil
}

func StopMjpgStreamer() error{
	process, err := os.FindProcess(common.GetMjpgPid())
	if err != nil {
		log.Printf("failed to find process by pid: %d, reason: %v\n", common.GetMjpgPid(), process)
		return err
	}
	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		log.Printf("failed to kill process %d: %v", common.GetMjpgPid(), err)
	} else {
		log.Println("killed process: ",  common.GetMjpgPid())
	}
	common.EditMjpgPid(-1)
	return nil
}