package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type Image struct {
	Src string `json:"src"`
}

// UploadImage 函数使用多协程上传图片
func UploadImage(cfPictureURL, filePath, fileFieldName string) error {
	postUrl := cfPictureURL + "/upload"

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("unable to open file: %v", err)
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("unable to get file info: %v", err)
	}

	// 计算块的数量
	chunkSize := 5 * 1024 * 1024 // 5MB
	numParts := int(math.Ceil(float64(fileInfo.Size()) / float64(chunkSize)))

	var wg sync.WaitGroup
	errChan := make(chan error, numParts)

	// 协程上传每个块
	for i := 0; i < numParts; i++ {
		wg.Add(1)
		go func(partNumber int) {
			defer wg.Done()
			partSize := chunkSize
			if i == numParts-1 {
				partSize = int(fileInfo.Size()) - (i * chunkSize)
			}

			// 创建一个新的缓冲区来读取每个部分
			filePart := make([]byte, partSize)

			// 定位到文件的当前块的起始位置
			_, err := file.Seek(int64(i)*int64(chunkSize), io.SeekStart)
			if err != nil {
				errChan <- fmt.Errorf("error seeking to file part %d: %v", i+1, err)
				return
			}

			// 读取当前块的数据
			_, err = io.ReadFull(file, filePart)
			if err != nil {
				errChan <- fmt.Errorf("error reading file part %d: %v", i+1, err)
				return
			}

			// 创建HTTP请求
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			part, err := writer.CreateFormFile(fileFieldName, fmt.Sprintf("%s.part%d", filepath.Base(file.Name()), partNumber))
			if err != nil {
				errChan <- err
				return
			}
			_, err = part.Write(filePart)
			if err != nil {
				errChan <- err
				return
			}
			writer.Close()

			req, err := http.NewRequest(http.MethodPost, postUrl, body)
			if err != nil {
				errChan <- err
				return
			}
			req.Header.Set("Content-Type", writer.FormDataContentType())

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				errChan <- err
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				errChan <- fmt.Errorf("failed to upload part %d, status: %s", partNumber, resp.Status)
			}

			responseBody := &bytes.Buffer{}
			_, err = responseBody.ReadFrom(resp.Body)
			if err != nil {
				errChan <- fmt.Errorf("error reading response body: %v", err)
				return
			}
			result := responseBody.String()

			var images []Image
			err = json.Unmarshal([]byte(result), &images)
			if err != nil {
				log.Fatalf("JSON parsing error: %v", err)
			}

			fmt.Println(cfPictureURL + images[0].Src)
		}(i)
	}

	// 等待所有协程完成
	wg.Wait()
	close(errChan)

	// 检查是否有错误发生
	if len(errChan) > 0 {
		for err := range errChan {
			fmt.Println("Error:", err)
		}
		return fmt.Errorf("one or more errors occurred during upload")
	}

	return nil
}
