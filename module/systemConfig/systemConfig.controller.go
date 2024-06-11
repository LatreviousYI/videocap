/*
 * @Author       : lvyitao lvyitao@fanhaninfo.com
 * @Date         : 2024-06-05 11:25:20
 * @LastEditTime : 2024-06-05 17:19:30
 */
package systemConfig

import (
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
	err = utils.Save(systemConfigModel.DataFileName(), systemConfigModel)
	if err != nil {
		utils.ErrResponse(c, 500, err.Error(), "")
		return
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
