/*
 * @Author       : lvyitao lvyitao@fanhaninfo.com
 * @Date         : 2024-06-07 16:41:03
 * @LastEditTime : 2024-06-07 16:45:48
 */
package systemConfig

import (

	"github.com/gin-gonic/gin"
)

func Init(factoryRouter *gin.Engine){
	systemConfigRouetr := factoryRouter.Group("/api/system/config")
	systemConfigRouetr.POST("/update",systemConfigUpdate)
	systemConfigRouetr.GET("/detail",systemConfigDetail)
	systemConfigRouetr.POST("/output/update",systemConfigOutputUpdate)
	systemConfigRouetr.GET("/output/detail",systemConfigOutputDetail)
}