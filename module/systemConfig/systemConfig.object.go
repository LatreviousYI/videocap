package systemConfig

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"src/common"
	"src/utils"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

var FileTemplateMap = map[string]string{"{{today}}":utils.Today(),"{{deviceid}}":DeviceId(),"{{datetime}}":utils.DateTime(),"{{unixtime}}":utils.UnixTime()}
var FileTemplateFuncMap = map[string]func()string{"{{today}}":utils.Today,"{{deviceid}}":DeviceId,"{{datetime}}":utils.DateTime,"{{unixtime}}":utils.UnixTime}


func(s *SystemConfigModel)UpdateVerify()error{
	if s.DeviceId == ""|| s.CollectionFrequency == ""|| s.ImgNameRules == ""{
		return errors.New("某个字段未填写")
	}
	s.CollectionFrequency = strings.TrimSpace(s.CollectionFrequency)
	specParser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	_, err := specParser.Parse(s.CollectionFrequency)
	if err != nil {
		log.Println(err.Error())
		return errors.New("采集时间点字段不合法")
	}
	if !utils.IsLegalString(s.DeviceId){
		return errors.New("设备编号不合法,只能包含中文字母数字或下划线") 
	}
	// 定义正则表达式
	re := regexp.MustCompile(`\{\{(.*?)\}\}`)

	// 查找所有匹配项
	matches := re.FindAllStringSubmatch(s.ImgNameRules, -1)
	if len(matches) == 0{
		return errors.New("图片保存命名规则不合法")
	}
	if s.ResolutionRatio != ""{
		ratioLen := strings.Split(s.ResolutionRatio, "x")
		if len(ratioLen) != 2{
			return errors.New("分辨率不合法")
		}
		ratioleftNum,err := strconv.Atoi(ratioLen[0])
		if err != nil || ratioleftNum < 200{
			return errors.New("分辨率数值不合法或太小")
		}
		ratiorightNum,err := strconv.Atoi(ratioLen[1])
		if err != nil || ratiorightNum < 200{
			return errors.New("分辨率数值不合法或太小")
		}
	}
	return nil
}

func(o *OutputConfig)UpdateVerify()error{
	if o.Local.Enable{
		if o.Local.OutputPath == ""{
			return errors.New("本地路径未填写")
		}
		outPutconfig := OutputConfig{}
		err :=utils.Get(outPutconfig.DataFileName(),&outPutconfig)
		if err != nil{
			return err
		}
		if outPutconfig.Local.OutputPath != o.Local.OutputPath{
			localPath := filepath.Join(o.Local.OutputPath,"test.txt")
			utils.PathExistsAndCreate(o.Local.OutputPath)
			f,err := os.Create(localPath)
			if err != nil{
				return err
			}
			defer f.Close()
			os.Remove(localPath)
			os.RemoveAll(o.Local.OutputPath)
			os.RemoveAll(outPutconfig.Local.OutputPath)
		}
	}
	if o.Clouds.Enable{
		if o.Clouds.Host == ""{
			return errors.New("host未填写")
		}
		o.Clouds.Host = strings.TrimRight(o.Clouds.Host,"/")
	}
	return nil
}

func ReplaceImgName(oldName string) string{
	for k,v:= range FileTemplateFuncMap{
		oldName = strings.ReplaceAll(oldName,k,v())
	}
	return oldName
}


func GetImage()([]byte,error){
	url := "http://127.0.0.1:8080/?action=snapshot"
	resp,err := http.Get(url)
	if err != nil{
		log.Println("http://127.0.0.1:8080/?action=snapshot获取图片报错",err.Error())
		return []byte{},err
	}
	if resp.StatusCode != 200{
		log.Println("http://127.0.0.1:8080/?action=snapshot获取图片状态码报错")
		return []byte{},err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil{
		log.Println(err)
		return []byte{},err
	}
    defer resp.Body.Close()
	return body,nil
}


func(s *SystemConfigModel)GetData()error{
	err := utils.Get(s.DataFileName(),s)
	if err != nil {
		return err
	}
	return nil
}

func ImageCapture(){
	outputConfig := OutputConfig{}
	err := utils.Get(outputConfig.DataFileName(),&outputConfig)
	if err != nil{
		log.Println(err)
		return
	}
	systemConfigModel := SystemConfigModel{}
	err = utils.Get(systemConfigModel.DataFileName(),&systemConfigModel)
	if err != nil{
		log.Println(err)
		return
	}
	body ,err := GetImage()
	if err != nil{
		log.Println(err)
		return
	}
	
	if outputConfig.Local.Enable{
		filePath := filepath.Join(outputConfig.Local.OutputPath,systemConfigModel.ImgNameRules)+".png"
		newfilpath :=  ReplaceImgName(filePath)
		go func(path string,content []byte){
			defer utils.ErrCatch()
			err := SaveLocalImage(path,content)
			if err != nil{
				log.Println(err)
			}
		}(newfilpath,body)
	}
	if outputConfig.Clouds.Enable{
		go SaveCloudsImage(body,systemConfigModel,outputConfig)
	}
}


func SaveLocalImage(path string,content []byte) error{
	dir,_ := filepath.Split(path)
	err := utils.PathExistsAndCreate(dir)
	if err != nil{
		return err
	}
	file,err:=os.Create(path)
	if err != nil{
		log.Println(err)
		return err
	}
	defer file.Close()
	_,err = file.Write(content)
	if err != nil{
		log.Println(err)
		return err
	}
	return nil
}

func SaveCloudsImage(imagebyte []byte,systemConfig SystemConfigModel,outputConfig OutputConfig)error{
	defer utils.ErrCatch()
	httpClient := &http.Client{}
	schemaSaveImageIn := SchemaSaveImageIn{}
	imagebs64 := base64.StdEncoding.EncodeToString(imagebyte)
	schemaSaveImageIn.Imagebs64 = "data:image/png;base64,"+imagebs64
	schemaSaveImageIn.DeviceID = systemConfig.DeviceId
	ip,err := utils.GetIP()
	if err != nil{
		log.Println(err)
		return err
	}
	schemaSaveImageIn.DeviceIP = ip
	schemaSaveImageIn.Datetime = utils.DateTime()
	requestBody, _ := json.Marshal(schemaSaveImageIn)
	req, _ := http.NewRequest("POST", outputConfig.Clouds.Host+"/api/v4/AI_httpService/DEV00085/SaveImage", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type","application/json")
	response, err := httpClient.Do(req)
	if err != nil {
		log.Println("保存图片服务请求失败")
		return err
	}
	body, _ := io.ReadAll(response.Body)
	defer response.Body.Close()
	if response.StatusCode > 201 && response.StatusCode < 299 {
		log.Println("保存图片服务请求失败StatusCode不等于200")
		return errors.New(response.Status)
	}
	var schemaSaveImageOut SchemaSaveImageOut
	err = json.Unmarshal(body, &schemaSaveImageOut)
	if err != nil{
		log.Println(string(body))
		log.Println(err.Error())
		return err
	}
	if !schemaSaveImageOut.Success {
		log.Println(string(body))
		log.Println("success不等于true")
		return errors.New("success不等于true")
	}
	return nil
}

func DeviceId() string{
	systemConfigModel := SystemConfigModel{}
	err := utils.Get(systemConfigModel.DataFileName(),&systemConfigModel)
	if err != nil{
		log.Println(err)
		return systemConfigModel.DeviceId
	}
	return systemConfigModel.DeviceId
}




func initDataFile(){
	defer utils.ErrCatch()
	systemConfigModel := SystemConfigModel{}
	err := utils.Get(systemConfigModel.DataFileName(),&systemConfigModel)
	if err != nil{
		systemConfigModel.CollectionFrequency = "0 0/5 * * * ?"
		systemConfigModel.ImgNameRules = "{{today}}/{{deviceid}}-{{datetime}}"
		systemConfigModel.ResolutionRatio = "640x480"
		machineId := utils.GetMd5(utils.GetMac())
		systemConfigModel.MachineId = machineId
		systemConfigModel.DeviceId = machineId
		err := utils.Save(systemConfigModel.DataFileName(),&systemConfigModel)
		if err != nil{
			log.Println("写入文件报错",err)
			return
		}
	}
	outputConfig := OutputConfig{}
	err = utils.Get(outputConfig.DataFileName(),&outputConfig)
	if err != nil{
		err := utils.Save(outputConfig.DataFileName(),&outputConfig)
		if err != nil{
			log.Println("写入文件报错",err)
			return
		}
	}
	entryid,err := common.CronTab.AddFunc(systemConfigModel.CollectionFrequency,ImageCapture)
	if err != nil{
		log.Println("创建定时任务报错",err)
		return
	}
	common.EditEntryId(entryid)
}


func GuaranteeMjpgServerRunning(){
	defer utils.ErrCatch()
	_,err := GetImage()
	if err != nil{
		var systemConfig SystemConfigModel
		err := systemConfig.GetData()
		if err != nil{
			log.Println(err.Error())
			return
		}
		var defaultDeviceList  = []string{"/dev/video1"}
		deList ,err := utils.GetDeviceList()
		if err != nil{
			log.Println(err.Error())
		}
		defaultDeviceList = append(defaultDeviceList, deList...)
		for _,v:= range defaultDeviceList{
			log.Println("尝试使用这个设备id启动MJPG-streamer:",v)
			err := utils.StartMjpgStreamer(systemConfig.ResolutionRatio,v)
			if err != nil{
				log.Println(err.Error())
			}
			time.Sleep(1*time.Second)
			_,err = GetImage()
			if err == nil{
				log.Println("启动MJPG-streamer成功")
				break
			}
		}
	}else{
		log.Println("MJPG-streamer已启动")
	}
}


func CleanImages(){
	defer utils.ErrCatch()
	outPutConfig := OutputConfig{}
	weekAgo := time.Now().Unix() -int64(604800)
	err := utils.Get(outPutConfig.DataFileName(),&outPutConfig)
	if err != nil{
		log.Println(err.Error())
		return
	}
	if outPutConfig.Local.Enable{
		fileList,err := os.ReadDir(outPutConfig.Local.OutputPath)
		if err != nil{
			log.Println(err.Error())
			return
		}
		for _,v:= range fileList{
			fileinfo,err:= v.Info()
			if err != nil{
				log.Println(err.Error())
				continue
			}
			fileModTime := fileinfo.ModTime().Unix()
			if weekAgo >= fileModTime{
				os.RemoveAll(filepath.Join(outPutConfig.Local.OutputPath,fileinfo.Name()))
			}
		}
	}
}


func TimeCleanlog(){
	defer utils.ErrCatch()
	CopyLog()
	CleanLog()
}
// cat /dev/null > aaa.log
func CleanLog(){
	basePath := utils.GetExcutePath()
	logfile := filepath.Join(basePath,"server.log")
	cmd := exec.Command("/bin/bash","-c",fmt.Sprintf("cat /dev/null > %s",logfile))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout  // 标准输出
	cmd.Stderr = &stderr  // 标准错误
	err := cmd.Run()
	if err != nil{
		log.Println(err.Error())
		return
	}
	_, errStr := stdout.String(), stderr.String()
	if errStr != ""{
		log.Println(errStr)
		return
	}
}


func CopyLog(){
	basePath := utils.GetExcutePath()
	logfile := filepath.Join(basePath,"server.log")
	cmd := exec.Command("/bin/bash","-c",fmt.Sprintf("cp %s %s",logfile,logfile+".bak"))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout  // 标准输出
	cmd.Stderr = &stderr  // 标准错误
	err := cmd.Run()
	if err != nil{
		log.Println(err.Error())
		return
	}
	_, errStr := stdout.String(), stderr.String()
	if errStr != ""{
		log.Println(errStr)
		return
	}
}


func(w *WifiConfig) WifiVerify()error{
	if w.Ssid =="" || w.Password == ""{
		return errors.New("wifi名称和密码未上传")
	}
	return nil
}