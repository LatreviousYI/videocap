/*
 * @Author       : lvyitao lvyitao@fanhaninfo.com
 * @Date         : 2024-06-07 15:56:38
 * @LastEditTime : 2024-06-11 11:02:51
 */
package cmd

import (
	"log"
	"os"
	"path/filepath"
	"src/utils"
	"strings"

	"github.com/spf13/cobra"
)

func init(){
	rootCmd.AddCommand(systemInstall)
}

func systemInstallfunc(){
	systemFile :=`[Unit]
Description=capVi Service

[Service]
Type=simple
ExecStart=#path# start console
KillSignal=SIGINT
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target`
	systemPath := "/usr/lib/systemd/system/capVi.service"
	basePath := utils.GetExcutePath()
	exePath := filepath.Join(basePath,"capVi")
	f,err := os.Create(systemPath)
	if err != nil{
		log.Println(err)
		return
	}
	defer f.Close()
	systemFile = strings.ReplaceAll(systemFile,"#path#",exePath)
	_,err = f.Write([]byte(systemFile))
	if err != nil{
		log.Println(err)
		return
	}
	err = utils.ReloadSystemctlDaemon()
	if err != nil{
		log.Println(err)
	}
}

var systemInstall = &cobra.Command{
	Use: "system_install",
	Short: "add systemctl server file",
	Run: func(cmd *cobra.Command, args []string) {
		systemInstallfunc()
	},
}


