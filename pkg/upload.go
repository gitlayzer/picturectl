package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type Image struct {
	Src string `json:"src"`
}

func UploadImage(cfPictureURL, filePath, fileFieldName string) error {
	// 拼接上传图片的URL
	postUrl := cfPictureURL + "/upload"
	// 读取文件
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("unable to open file: %v", err)
	}
	defer file.Close()

	// 创建一个缓冲区来存储 multipart 的表单数据
	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)

	// 创建表单文件部分
	formFile, err := multipartWriter.CreateFormFile(fileFieldName, filepath.Base(file.Name()))
	if err != nil {
		return fmt.Errorf("partial failure to create form file: %v", err)
	}

	// 将文件数据复制到表单
	_, err = io.Copy(formFile, file)
	if err != nil {
		return fmt.Errorf("failed to copy file to form: %v", err)
	}

	// 关闭multipart writer来设置表单的最终边界
	multipartWriter.Close()

	// 创建HTTP请求
	request, err := http.NewRequest(http.MethodPost, postUrl, &requestBody)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}
	// 设置Content-Type头部，这样服务器知道是multipart/form-data
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("sending request failed: %v", err)
	}
	defer response.Body.Close()

	// 检查响应状态码
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %v", response.Status)
	}

	// 打印响应内容
	responseBody := &bytes.Buffer{}
	_, err = responseBody.ReadFrom(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}
	result := responseBody.String()

	// 解析JSON
	var images []Image
	err = json.Unmarshal([]byte(result), &images)
	if err != nil {
		log.Fatalf("JSON parsing error: %v", err)
	}

	// 检查是否有结果并打印出来
	fmt.Println(cfPictureURL + images[0].Src)

	return nil
}
