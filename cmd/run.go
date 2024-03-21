package cmd

import (
	"fmt"
	"github.com/gitlayzer/picturectl/pkg"
	"log"
	"os"
)

func Run() {
	// 判断第一个参数是否为一个URL
	if len(os.Args) < 3 {
		fmt.Println("Usage: picturectl <pictureUrl> <imagePath>")
		return
	}

	pictureUrl := os.Args[1]

	imagePath := os.Args[2]

	err := pkg.UploadImage(pictureUrl, imagePath, "file")
	if err != nil {
		log.Fatalf("Failed to upload image: %v", err)
	}
}
