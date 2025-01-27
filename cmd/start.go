package cmd

import (
	"context"
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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func init(){
	rootCmd.AddCommand(startCmd)
}

func start(){
	go func() {
		// 服务连接
		factory.Init()
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
	Short: "Start the application    if start console 启动时会打开控制台,不会再后台运行",
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 && args[0] == "console"{
			start()
		}else{
			daemonStart()
		}
	},
}

func ServerExit() {
	log.Println("开始退出")
	factory.GlobalCancel()
	factory.GlobalWG.Wait()
	// 清理pid文件
	err := os.Remove(pidFile)
	if err != nil {
		log.Println("failed to remove pid file",err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := factory.HttpServer.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	os.Exit(0)
}