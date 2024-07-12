'''
Author       : kerwin_lv 2509456238@qq.com
Date         : 2024-01-31 14:14:54
LastEditTime : 2024-04-16 16:24:10
'''
import json
import os,shutil,pathlib,subprocess,argparse

root_path = pathlib.Path(os.path.abspath(__file__)).parent
project_path = root_path.parent.parent

# 创建发布文件夹
publish_path = os.path.join(root_path,"publish")
if not os.path.exists(publish_path):
    os.makedirs(publish_path)
else:
    shutil.rmtree(publish_path)
    os.makedirs(publish_path)
    

# 需要复制
config_path = os.path.join(root_path,"config")
uiptah_path = os.path.join(root_path,"mjpg-streamer")
wwwroot_path = os.path.join(root_path,"www")
exe_path = os.path.join(root_path,"capVi")
# 需要创建
db_path = os.path.join(publish_path,"data")

publish_config_path = os.path.join(publish_path,"config")
publish_uipath_path = os.path.join(publish_path,"mjpg-streamer")
publish_wwwroot_path = os.path.join(publish_path,"www")


result = subprocess.run('go build -o capVi ./main.go',shell=True,cwd=root_path,capture_output=True,text=True)
print(result.stderr)
print("build")

shutil.copytree(config_path, publish_config_path)
shutil.copytree(uiptah_path, publish_uipath_path)
shutil.copytree(wwwroot_path, publish_wwwroot_path)
shutil.copy(exe_path, publish_path)
os.makedirs(db_path)
shutil.make_archive("capVi", 'zip', publish_path)
print("发布成功")


