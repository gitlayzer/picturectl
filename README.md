# Picturectl Typora Upload Image Plugin

## 0. Overview
这个小工具是根据 [Telegraph-Image](https://github.com/cf-pages/Telegraph-Image) 这个图床项目实现的 Typora 的插件

## 1. Introduction
确保自己本地有 Go 语言的开发环境
- 代码克隆 git clone https://github.com/gitlayzer/picturectl.git
- 编译 windows 平台：make build-windows
- 编译 linux 平台：make build-linux
- 生成的文件在 bin 目录下，可以根据自己的需求放到合适的位置使用

## 2. Usage
这里以 windows 作为演示
- make build-windows
- 将生成的 picturectl.exe 放到你的 Typora 图片上传插件目录下
- 打开 Typora 设置，找到偏好设置
- 找到图像，在插入图片时选择上传图片
- 在上传服务设定的地方选择自定义命令
- 在命令内写入自己编译生成的 picturectl.exe 的路径
- 需要注意的是，这里需要在命令后面跟上你的图片的服务器地址
- 比如：picturectl.exe F:\Software\Typora\picturectl.exe https://picture.xxxxxxxxxxx.com.cn
![image-20240321154930617](https://picture.devops-engineer.com.cn/file/0609089c525fe6b0dad1e.png)