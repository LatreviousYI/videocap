package systemConfig

import (
	"errors"
	"log"
	"regexp"
	"src/utils"
	"strings"

	"github.com/robfig/cron/v3"
)



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
	return nil
}

func(o *OutputConfig)UpdateVerify()error{
	if o.Local.Enable{
		if o.Local.OutputPath == ""{
			return errors.New("本地路径未填写")
		}
	}
	if o.Nfs.Enable{
		if o.Nfs.OutputPath == "" || o.Nfs.NfsHost == ""{
			return errors.New("nfs未填写")
		}
	}
	if o.Cifs.Enable{
		if o.Cifs.OutputPath == ""||o.Cifs.CifsHost == ""||o.Cifs.Username == ""||o.Cifs.Password == ""{
			return errors.New("cifs未填写")
		}
	}
	return nil
}