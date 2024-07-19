/*
 * @Author       : lvyitao lvyitao@fanhaninfo.com
 * @Date         : 2024-06-07 16:41:03
 * @LastEditTime : 2024-06-07 16:45:48
 */
package systemConfig

import (
	"log"
	"src/common"

	"github.com/gin-gonic/gin"
)

func Init(factoryRouter *gin.Engine){
	systemConfigRouetr := factoryRouter.Group("/api/system/config")
	systemConfigRouetr.POST("/update",systemConfigUpdate)
	systemConfigRouetr.GET("/detail",systemConfigDetail)
	systemConfigRouetr.POST("/output/update",systemConfigOutputUpdate)
	systemConfigRouetr.GET("/output/detail",systemConfigOutputDetail)

	systemConfigRouetr.GET("/images",getImageContent)

	systemConfigRouetr.GET("/wifi/list",getWifiList)
	systemConfigRouetr.POST("/wifi/connect",connectWifi)
	systemConfigRouetr.GET("/wifi/status",checkWifiStatus)

	systemConfigRouetr.GET("/name/rules",getNameRulesList)
	
	initDataFile()
	GuaranteeMjpgServerRunning()

	_,err := common.CronTab.AddFunc("0 0 9 ? * MON",CleanImages)
	if err != nil{
		log.Println(err.Error())
		return
	}
	_,err = common.CronTab.AddFunc("0 0 9 ? * MON",TimeCleanlog)
	if err != nil{
		log.Println(err.Error())
		return
	}
}