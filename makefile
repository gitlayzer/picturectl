# 编译 windows 环境下的 Go 程序
build-windows:
	go build -o bin/picturectl.exe main.go

# 编译 linux 环境下的 Go 程序
build-linux:
	go build -o bin/picturectl main.go

# 编译 macOS 环境下的 Go 程序
build-macos:
	go build -o bin/picturectl main.go

# 运行本地的 Go 程序
run:
	go run main.go