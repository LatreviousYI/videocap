/*
 * @Author       : lvyitao lvyitao@fanhaninfo.com
 * @Date         : 2024-06-05 11:25:20
 * @LastEditTime : 2024-06-05 17:19:30
 */
package systemConfig

import (
	"src/common"
	"src/utils"

	"github.com/gin-gonic/gin"
)

func systemConfigUpdate(c *gin.Context) {
	systemConfigModel := SystemConfigModel{}
	if err := c.ShouldBind(&systemConfigModel); err != nil {
		utils.ErrResponse(c, 500, err.Error(), "")
		return
	}
	err := systemConfigModel.UpdateVerify()
	if err != nil {
		utils.ErrResponse(c, 500, err.Error(), "")
		return
	}
	oldSystemConfigModel := SystemConfigModel{}
	err = utils.Get(oldSystemConfigModel.DataFileName(),&oldSystemConfigModel)
	if err != nil {
		utils.ErrResponse(c, 500, err.Error(), "")
		return
	}
	systemConfigModel.MachineId = oldSystemConfigModel.MachineId
	err = utils.Save(systemConfigModel.DataFileName(), systemConfigModel)
	if err != nil {
		utils.ErrResponse(c, 500, err.Error(), "")
		return
	}
	if oldSystemConfigModel.CollectionFrequency != systemConfigModel.CollectionFrequency{
		common.CronTab.Remove(common.GetEntryId())
		entryid,err := common.CronTab.AddFunc(systemConfigModel.CollectionFrequency,ImageCapture)
		if err != nil{
			utils.ErrResponse(c, 500, err.Error(), "")
			return
		}
		common.EditEntryId(entryid)
	}
	if oldSystemConfigModel.ResolutionRatio != systemConfigModel.ResolutionRatio{
		utils.StopMjpgStreamer()
		GuaranteeMjpgServerRunning()
	}
	utils.SuccessResp(c, "")
}

func systemConfigDetail(c *gin.Context) {
	systemConfigModel := SystemConfigModel{}
	err := utils.Get(systemConfigModel.DataFileName(), &systemConfigModel)
	if err != nil {
		utils.ErrResponse(c, 500, err.Error(), "")
		return
	}
	utils.SuccessResp(c, systemConfigModel)
}

func systemConfigOutputUpdate(c *gin.Context) {
	outputConfig := OutputConfig{}
	if err := c.ShouldBind(&outputConfig); err != nil {
		utils.ErrResponse(c, 500, err.Error(), "")
		return
	}
	err := outputConfig.UpdateVerify()
	if err != nil {
		utils.ErrResponse(c, 500, err.Error(), "")
		return
	}
	err = utils.Save(outputConfig.DataFileName(), outputConfig)
	if err != nil {
		utils.ErrResponse(c, 500, err.Error(), "")
		return
	}
	utils.SuccessResp(c, "")

}

func systemConfigOutputDetail(c *gin.Context) {
	outputConfig := OutputConfig{}
	err := utils.Get(outputConfig.DataFileName(), &outputConfig)
	if err != nil {
		utils.ErrResponse(c, 500, err.Error(), "")
		return
	}
	utils.SuccessResp(c, outputConfig)
}

func getImageContent(c *gin.Context){
	b,err:=GetImage()
	if err != nil{
		utils.ErrResponse(c, 500, err.Error(), "")
		return
	}
	c.Data(200,"image/png",b)
}

func getWifiList(c *gin.Context){
	wifiList,err:=utils.GetWifiList()
	if err != nil{
		utils.ErrResponse(c, 500, err.Error(), "")
		return
	}
	utils.SuccessResp(c,wifiList)
}

func connectWifi(c *gin.Context){
	wifiInfo := WifiConfig{}
	if err := c.ShouldBind(&wifiInfo); err != nil {
		utils.ErrResponse(c, 500, err.Error(), "")
		return
	}
	err := utils.ConnectWifi(wifiInfo.Ssid,wifiInfo.Password)
	if err != nil{
		utils.ErrResponse(c, 500, err.Error(), "")
		return
	}
	utils.SuccessResp(c,"")
}

func checkWifiStatus(c *gin.Context){
	ok,err:=utils.CheckWifiStatus()
	if err != nil{
		utils.ErrResponse(c, 500, err.Error(), "")
		return
	}
	if !ok{
		utils.ErrResponse(c, 600, "未连接", "")
		return
	}
	utils.SuccessResp(c,ok)
}

func getNameRulesList(c *gin.Context){
	utils.SuccessResp(c,FileTemplateMap)
}