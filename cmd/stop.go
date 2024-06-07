/*
 * @Author       : lvyitao lvyitao@fanhaninfo.com
 * @Date         : 2024-06-07 15:56:38
 * @LastEditTime : 2024-06-07 16:10:03
 */
package cmd

import (
	"context"
	"log"
	"os"
	"src/factory"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

func init(){
	rootCmd.AddCommand(stopCmd)
}

func stop(){
	initDaemon()
	if pid == -1 {
		log.Println("Seems not have been started. Try use `alist start` to start server.")
		return
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		log.Printf("failed to find process by pid: %d, reason: %v\n", pid, process)
		return
	}
	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		log.Printf("failed to kill process %d: %v", pid, err)
	} else {
		log.Println("killed process: ", pid)
	}
	pid = -1
}

var stopCmd = &cobra.Command{
	Use: "stop",
	Short: "stop the application",
	Run: func(cmd *cobra.Command, args []string) {
		stop()
	},
}


func ServerExit() {
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