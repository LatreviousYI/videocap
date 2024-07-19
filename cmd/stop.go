/*
 * @Author       : lvyitao lvyitao@fanhaninfo.com
 * @Date         : 2024-06-07 15:56:38
 * @LastEditTime : 2024-06-11 11:02:51
 */
package cmd

import (
	"log"
	"os"
	"syscall"
	"github.com/spf13/cobra"
)

func init(){
	rootCmd.AddCommand(stopCmd)
}

func stop(){
	initDaemon()
	if pid == -1 {
		log.Println("Seems not have been started. Try use `./capVi start` to start server.")
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
		log.Println("stop success killed process: ", pid)
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


