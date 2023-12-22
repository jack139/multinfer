package ocr

import (
	"fmt"
	"log"
	"strconv"

	"github.com/jack139/go-infer/helper"
)



/*  定义模型相关参数和方法  */
type OCRIdCard struct{}

func (x *OCRIdCard) Init() error {
	return nil
}

func (x *OCRIdCard) ApiPath() string {
	return "/ocr2/id_card"
}

func (x *OCRIdCard) CustomQueue() string {
	return helper.Settings.Customer["OCR_QUEUE"]
}

func (x *OCRIdCard) ApiEntry(reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Api_OCRIdCard")

	// 检查参数
	imageBase64, ok := (*reqData)["image"].(string)
	if !ok {
		return &map[string]interface{}{"code":9101}, fmt.Errorf("need image")
	}

	// 检查图片大小
	maxSize, _ := strconv.Atoi(helper.Settings.Customer["OCR_MAX_IMAGE_SIZE"])
	if len(imageBase64) > maxSize {
		return &map[string]interface{}{"code":9002}, fmt.Errorf("图片数据太大")
	}

	// 构建请求参数
	reqDataMap := map[string]interface{}{
		"image": imageBase64,
	}

	return &reqDataMap, nil
}


// OCRIdCard 推理 - 不在这里实现，由 python dispatcher 实现
func (x *OCRIdCard) Infer(reqId string, reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Infer_OCRIdCard - Do nothing")

	return &map[string]interface{}{}, nil
}
