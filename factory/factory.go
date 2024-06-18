/*
 * @Author       : lvyitao lvyitao@fanhaninfo.com
 * @Date         : 2024-06-05 11:20:39
 * @LastEditTime : 2024-06-07 16:46:17
 */
package factory

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"src/module/systemConfig"
	"src/utils"
	"sync"

	"github.com/gin-gonic/gin"
)

var GinRouter *gin.Engine

var HttpServer *http.Server

var GlobalWG sync.WaitGroup

var GlobalCtx context.Context

var GlobalCancel context.CancelFunc

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile | log.Ldate)
	conifg := utils.GetConfig()
	if !conifg.Debug{
		gin.SetMode(gin.ReleaseMode)
	}
	GinRouter = gin.Default()
	GinRouter.Use(Cors())
	GinRouter.Use(Recover)
	HttpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", conifg.Port),
		Handler: GinRouter,
	}
	GlobalCtx, GlobalCancel = context.WithCancel(context.Background())

	systemConfig.Init(GinRouter)
	
}


