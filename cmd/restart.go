/*
 * @Author       : lvyitao lvyitao@fanhaninfo.com
 * @Date         : 2024-06-07 16:22:22
 * @LastEditTime : 2024-06-07 16:24:54
 */
package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

func init(){
	rootCmd.AddCommand(restartCommand)
}


var restartCommand = &cobra.Command{
	Use: "restart",
	Short: "Restart the application",
	Run: func(cmd *cobra.Command, args []string) {
		stop()
		time.Sleep(1*time.Second)
		daemonStart()
	},
}