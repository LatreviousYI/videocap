package cmd

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"src/factory"
	"src/utils"
	"strconv"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func init(){
	rootCmd.AddCommand(startCmd)
}

func start(){
	factory.GinRouter.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome Gin Server")
	})
	go func() {
		// 服务连接
		log.Println("server start")
		if err := factory.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	ServerExit()
	log.Println("Server exiting")
}


func daemonStart(){
	initDaemon()
	if pid != -1{
		_, err := os.FindProcess(pid)
		if err == nil {
			log.Println("alist already started, pid ", pid)
			return
		}
	}
	args := os.Args
	args[1] = "start"
	args = append(args, "console")
	cmd := &exec.Cmd{
		Path: args[0],
		Args: args,
		Env:  os.Environ(),
	}
	gin.SetMode(gin.ReleaseMode)
	basePath := utils.GetExcutePath()
	stdout, err := os.OpenFile(filepath.Join(basePath, "server.log"), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(os.Getpid(), ": failed to open start log file:", err)
	}
	cmd.Stderr = stdout
	cmd.Stdout = stdout
	err = cmd.Start()
	if err != nil {
		log.Fatal("failed to start children process: ", err)
	}
	log.Printf("success start pid: %d\n", cmd.Process.Pid)
	err = os.WriteFile(pidFile, []byte(strconv.Itoa(cmd.Process.Pid)), 0666)
	if err != nil {
		log.Println("failed to record pid, you may not be able to stop the program with `./main stop`")
	}
	log.Println("server start")
}



var startCmd = &cobra.Command{
	Use: "start",
	Short: "Start the application",
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 && args[0] == "console"{
			start()
		}else{
			daemonStart()
		}
	},
}

