package yhface

import (
	"os"
	"time"
	"fmt"
	"log"
	"math"
	"io/ioutil"

	"github.com/jack139/go-infer/helper"
)

// 保存请求图片和结果
func saveBackLog(requestId string, image, result []byte) {
	if helper.Settings.Customer["FACE_SAVE_IMAGE"] == "1" {
		output_dir := fmt.Sprintf("%s/%s", 
			helper.Settings.Customer["FACE_SAVE_IMAGE_PATH"], 
			time.Now().Format("20060102"))
		err := os.Mkdir(output_dir, 0755) // 建日志目录， 日期 做子目录
		if err == nil || os.IsExist(err) { // 不处理错误
			_ = ioutil.WriteFile(fmt.Sprintf("%s/%s.jpg", output_dir, requestId), image, 0644)
			_ = ioutil.WriteFile(fmt.Sprintf("%s/%s.txt", output_dir, requestId), result, 0644)
		} else {
			log.Println("ERROR when saving log: ", err.Error())
		}
	}
}

// 向量余弦相似度
func cosine(a []float32, b []float32) (cosine float64, err error) {
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
