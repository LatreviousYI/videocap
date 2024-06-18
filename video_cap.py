'''
Author       : lvyitao lvyitao@fanhaninfo.com
Date         : 2024-06-12 10:50:10
LastEditTime : 2024-06-12 10:58:29
'''
import cv2

# 打开摄像头
cap = cv2.VideoCapture(1)

ret, frame = cap.read()

# if not ret:
print(cap.isOpened())
# 显示视频帧
cv2.imwrite("l.png",frame)
# cv2.imshow('Video', frame)

# # 按下ESC键退出
# if cv2.waitKey(1) == 27:
#     break
    

# 释放摄像头并关闭窗口
cap.release()
cv2.destroyAllWindows()