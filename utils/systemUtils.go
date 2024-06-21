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
	if err != nil {
		return wifiList,err
	}
	outStr, errStr := stdout.String(), stderr.String()
	if errStr != ""{
		return wifiList,errors.New(errStr)
	}
	rowWifiList := strings.Split(strings.ReplaceAll(outStr,"\r\n", "\n"),"\n")
	for i,v:= range rowWifiList{
		if i == 0{
			continue
		}
		v = strings.ReplaceAll(v,"*","")
		oneWifiList := strings.Fields(v)
		if len(oneWifiList) < 3{
			continue
		}
		var wifi = Wifi{
			Ssid: oneWifiList[1],
			RowInfo: strings.TrimLeft(v," "),
		}
		wifiList = append(wifiList, wifi)
	}
	return wifiList,nil
}

func CheckWifiStatus()(bool,error){
	cmd := exec.Command("/bin/bash","-c",string("nmcli dev status"))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout  // 标准输出
	cmd.Stderr = &stderr  // 标准错误
	err := cmd.Run()
	if err != nil {
		return false,err
	}
	outStr, errStr := stdout.String(), stderr.String()
	if errStr != ""{
		return false,errors.New(errStr)
	}
	rowWifistatusList := strings.Split(strings.ReplaceAll(outStr,"\r\n", "\n"),"\n")
	for i,v:= range rowWifistatusList{
		if i == 0{
			continue
		}
		vList := strings.Fields(v)
		if vList[2] == "connected"{
			return true,nil
		}
	}
	return false,nil
}

func ConnectWifi(wifiname,password string) error{
	cmd := exec.Command("/bin/bash","-c",fmt.Sprintf("nmcli dev wifi connect %s password %s",wifiname,password))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout  // 标准输出
	cmd.Stderr = &stderr  // 标准错误
	err := cmd.Run()
	if err != nil {
		return err
	}
	outStr, errStr := stdout.String(), stderr.String()
	if errStr != ""{
		return errors.New(errStr)
	}
	if !strings.Contains(outStr,"successfully"){
		return errors.New(outStr)
	}
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
