'''
Author       : lvyitao lvyitao@fanhaninfo.com
Date         : 2024-05-27 09:47:07
LastEditTime : 2024-07-03 10:03:40
'''
import argparse
import os
import pathlib
import requests

root_path = pathlib.Path(os.path.abspath(__file__)).parent
project_path = root_path.parent.parent

project_id = 90
host_port = "10.21.77.101:9007"
upload_file_url = f"http://{host_port}/api/v1/createPublishVersion"
file_path = os.path.join(root_path,"capVi.zip")

parser = argparse.ArgumentParser(description='创建发布中心版本脚本')
# 给这个解析对象添加命令行参数
parser.add_argument('version', type=str, help='set product version')

args = parser.parse_args()
if args.version == "":
    print("未接收版本号 比如:python release.py 2.253")
    os._exit(0)

print("version:",args.version)

publish_data = {
    "version":args.version,
    "project_id":project_id,
    "remark":"",
    "is_upload_cloud":"1"
}

with open(file_path, 'rb') as f:
    upload_response = requests.post(
        url=upload_file_url,
        data= publish_data,
        files={"file": ("capVi.zip",f,"application/zip")}
    )
    print(upload_response.text)

# 下载链接
# 内网
# http://10.21.77.101:9007/api/getPublishFileUseProjectName?project_name=capVi
# 外网
# http://103.222.40.154:9007/api/getPublishFileUseProjectName?project_name=capVi&is_upload_cloud=1
