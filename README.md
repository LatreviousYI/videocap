# video-capture-terminal-server
```
capVi start 后台运行
capVi start console 终端等待运行
capVi stop 停止
capVi restart 重启
```

### 下载地址
```
内网
http://10.21.77.101:9007/api/getPublishFileUseProjectName?project_name=capVi
公网
http://103.222.40.154:9007/api/getPublishFileUseProjectName?project_name=capVi&is_upload_cloud=1
```

## 热点管理工具
### create_ap下载安装
```
git clone https://github.com/oblique/create_ap.git
cd create_ap
sudo make install
```
### create_ap conf配置
```
SHARE_METHOD=nat
WIFI_IFACE=wlan0
INTERNET_IFACE=eth0
SSID=orangepi
PASSPHRASE=orangepi
NO_VIRT=1
DAEMONIZE=0
```

### create_ap修复wifi无法连接
```
create_ap --fix-unmanaged
```


