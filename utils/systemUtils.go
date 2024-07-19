package utils

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"syscall"

	"github.com/shirou/gopsutil/process"
)


type Wifi struct{
	Ssid string `json:"ssid"`
	RowInfo string `json:"row_info"`
}


func GetWifiList() ([]Wifi,error){
	var wifiList []Wifi
	cmd := exec.Command("/bin/bash","-c",string("nmcli dev wifi"))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout  // 标准输出
	cmd.Stderr = &stderr  // 标准错误
	err := cmd.Run()
	outStr, errStr := stdout.String(), stderr.String()
	if errStr != "" || err != nil{
		log.Println(err)
		return wifiList,errors.New(errStr)
	}
	rowWifiList := strings.Split(strings.ReplaceAll(outStr,"\r\n", "\n"),"\n")
	for i,v:= range rowWifiList{
		if i == 0{
			continue
		}
		oneWifiList := strings.Fields(v)
		if len(oneWifiList) < 3{
			continue
		}
		var wifi Wifi
		if oneWifiList[0] == "*"{
			wifi = Wifi{
				Ssid: oneWifiList[2],
				RowInfo: strings.TrimLeft(v," "),
			}
		}else{
			wifi = Wifi{
				Ssid: oneWifiList[1],
				RowInfo: strings.TrimLeft(v," "),
			}
		}
		wifiList = append(wifiList, wifi)
	}
	return wifiList,nil
}

func CheckWifiStatus()(bool,error){
	cmd := exec.Command("/bin/bash","-c",string("nmcli network connectivity"))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout  // 标准输出
	cmd.Stderr = &stderr  // 标准错误
	err := cmd.Run()
	outStr, errStr := stdout.String(), stderr.String()
	if errStr != "" || err != nil{
		log.Println(err)
		return false,errors.New(errStr)
	}
	if outStr == ""{
		return false,errors.New("输出为空")
	}
	if !strings.Contains(outStr,"full"){
		return false,errors.New(outStr)
	}
	return true,nil
}

func StopCreateAp()error{
	cmd := exec.Command("/bin/bash","-c","systemctl stop create_ap.service")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout  // 标准输出
	cmd.Stderr = &stderr  // 标准错误
	err := cmd.Run()
	_, errStr := stdout.String(), stderr.String()
	if errStr != "" || err != nil{
		log.Println(errStr,err)
		return errors.New(errStr)
	}
	return nil
}


func ReloadSystemctlDaemon()error{
	cmd := exec.Command("/bin/bash","-c","systemctl daemon-reload")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout  // 标准输出
	cmd.Stderr = &stderr  // 标准错误
	err := cmd.Run()
	_, errStr := stdout.String(), stderr.String()
	if errStr != "" || err != nil{
		log.Println(errStr,err)
		return errors.New(errStr)
	}
	return nil
}


func FixWifi() error{
	cmd := exec.Command("/bin/bash","-c","create_ap --fix-unmanaged")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout  // 标准输出
	cmd.Stderr = &stderr  // 标准错误
	err := cmd.Run()
	_, errStr := stdout.String(), stderr.String()
	if errStr != "" || err != nil{
		log.Println(errStr,err)
		return errors.New(errStr)
	}
	return nil
}


func CloseUseCreateAp() error{
	cmd := exec.Command("/bin/bash","-c","systemctl disable create_ap.service")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout  // 标准输出
	cmd.Stderr = &stderr  // 标准错误
	err := cmd.Run()
	_, errStr := stdout.String(), stderr.String()
	if errStr != "" || err != nil{
		log.Println(errStr,err)
		return errors.New(errStr)
	}
	return nil
}

// create_ap --fix-unmanaged
func ConnectWifi(wifiname,password string) error{
	// ExecEnv()
	cmd := exec.Command("/bin/bash","-c",fmt.Sprintf(`nmcli device wifi connect %s password %s`,wifiname,password))
	// cmd := exec.Command("/bin/nmcli","device","wifi","connect",strconv.Quote(wifiname) ,"password",strconv.Quote(password))
	// cmd := exec.Command("./wifi_connect.sh",wifiname,password)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout  // 标准输出
	cmd.Stderr = &stderr  // 标准错误
	err := cmd.Run()
	outStr, errStr := stdout.String(), stderr.String()
	if errStr != "" || err != nil{
		log.Println(errStr,err)
		return errors.New(errStr)
	}
	if !strings.Contains(outStr,"successfully"){
		log.Println(outStr)
		return errors.New(outStr)
	}
	return nil
}


func ExecEnv()error{
	cmd := exec.Command("env")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout  // 标准输出
	cmd.Stderr = &stderr  // 标准错误
	err := cmd.Run()
	if err != nil {
		log.Println(err)
		return err
	}
	outStr, errStr := stdout.String(), stderr.String()
	log.Println("out",outStr)
	log.Println("err",errStr)
	return nil
}

func GetMac() string {
	goos := runtime.GOOS
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("fail to get net interfaces: %v\n", err)
	}
	if len(netInterfaces) == 0 {
		log.Println(err.Error())
	}
	sort.Slice(netInterfaces, func(i, j int) bool {
		return netInterfaces[i].Index < netInterfaces[j].Index
	})
	switch goos {
	case "windows":
		for _, netInterface := range netInterfaces {
			macAddr := netInterface.HardwareAddr.String()
			if len(macAddr) == 0 {
				continue
			}
			if strings.HasPrefix(strings.ToLower(netInterface.Name), "以太网") {
				return macAddr
			}
		}
		for _, netInterface := range netInterfaces {
			macAddr := netInterface.HardwareAddr.String()
			if len(macAddr) == 0 {
				continue
			}

			if strings.HasPrefix(strings.ToLower(netInterface.Name), "eth") {
				return macAddr
			}
		}
		for _, netInterface := range netInterfaces {
			macAddr := netInterface.HardwareAddr.String()
			if len(macAddr) == 0 {
				continue
			}
			if strings.HasPrefix(strings.ToLower(netInterface.Name), "wl") {
				return macAddr
			}
		}

	case "linux":
		for _, netInterface := range netInterfaces {
			macAddr := netInterface.HardwareAddr.String()
			if len(macAddr) == 0 {
				continue
			}
			if strings.HasPrefix(strings.ToLower(netInterface.Name), "en") {
				return macAddr
			}
		}
		for _, netInterface := range netInterfaces {
			macAddr := netInterface.HardwareAddr.String()
			if len(macAddr) == 0 {
				continue
			}
			if strings.HasPrefix(strings.ToLower(netInterface.Name), "et") {
				return macAddr
			}
		}
		for _, netInterface := range netInterfaces {
			macAddr := netInterface.HardwareAddr.String()
			if len(macAddr) == 0 {
				continue
			}
			if strings.HasPrefix(strings.ToLower(netInterface.Name), "wla") {
				return macAddr
			}
		}
	default:
		for _, netInterface := range netInterfaces {
			macAddr := netInterface.HardwareAddr.String()
			if len(macAddr) == 0 {
				continue
			}
			return macAddr
		}
	}
	return ""
}


func GetIP() (string,error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "",err
	}
	defer conn.Close()
	ipAddress := conn.LocalAddr().(*net.UDPAddr)
	return ipAddress.IP.String(),nil
}


func GetProcessName(pid int)(string,error){
	proc, err := process.NewProcess(int32(pid))
    if err != nil {
		return "",err
    }
    // Get the process name
    name, err := proc.Name()
    if err != nil {
		return "",err
    }
	return name,nil
}

func IsProcessAlive(pid int) bool {
    // syscall.Kill with signal 0 does not actually send a signal,
    // but it checks for the existence of the process.
    err := syscall.Kill(pid, 0)
    return err == nil
}