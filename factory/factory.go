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
	"path/filepath"
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


func index(c *gin.Context){
	c.Header("content-type", "text/html;charset=utf-8")
	basePath := utils.GetExcutePath()
	c.File(filepath.Join(basePath,"www","index.html"))
}

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
	GinRouter.GET("/healthy", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome Gin Server")
	})
	GinRouter.Static("/static","./www/static")
	GinRouter.NoRoute(index)
	systemConfig.Init(GinRouter)
}


