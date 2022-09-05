package yhface

import (
	"os"
	"time"
	"fmt"
	"log"
	"math"
	"image"
	"bytes"
	"io/ioutil"
	"github.com/disintegration/imaging"
	"github.com/jack139/go-infer/helper"
)

// 保存请求图片和结果
func saveBackLog(requestId string, img image.Image, result []byte) {
	if helper.Settings.Customer["FACE_SAVE_IMAGE"] == "1" {
		output_dir := fmt.Sprintf("%s/%s", 
			helper.Settings.Customer["FACE_SAVE_IMAGE_PATH"], 
			time.Now().Format("20060102"))
		err := os.Mkdir(output_dir, 0755) // 建日志目录， 日期 做子目录
		if err == nil || os.IsExist(err) { // 不处理错误
			_ = imaging.Save(img, fmt.Sprintf("%s/%s.jpg", output_dir, requestId))
			_ = ioutil.WriteFile(fmt.Sprintf("%s/%s.txt", output_dir, requestId), result, 0644)
		} else {
			log.Println("ERROR when saving log: ", err.Error())
		}
	}
}

// 向量余弦相似度
func cosine(a []float32, b []float32) (float64, error) {
	sumA := 0.0
	s1 := 0.0
	s2 := 0.0

	for k := 0; k < len(a); k++ {
		sumA += float64(a[k]) * float64(b[k])
		s1 += math.Pow(float64(a[k]), 2)
		s2 += math.Pow(float64(b[k]), 2)
	}
	if s1 == 0 || s2 == 0 {
		return 0.0, fmt.Errorf("Vectors should not be null (all zeros)")
	}
	return - sumA / (math.Sqrt(s1) * math.Sqrt(s2)), nil
}

// 向量正则化
func norm(a []float32) ([]float32, error) {
	s1 := 0.0

	for k := 0; k < len(a); k++ {
		s1 += math.Pow(float64(a[k]), 2)
	}
	if s1 == 0 {
		return nil, fmt.Errorf("Vectors should not be null (all zeros)")
	}
	norm :=  float32(math.Sqrt(s1))

	for k := 0; k < len(a); k++ {
		a[k] = a[k] / norm
	}

	return a, nil
}


// image 转换为 bytes
func image2bytes(img image.Image) ([]byte, error) {
	// 转换为 []byte
	buf := new(bytes.Buffer)
	err := imaging.Encode(buf, img, imaging.JPEG)
	if err!=nil {
		return nil, err
	}

	return buf.Bytes(), nil
}