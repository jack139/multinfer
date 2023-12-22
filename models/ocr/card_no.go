package ocr

import (
	"fmt"
	"log"
	"strconv"

	"github.com/jack139/go-infer/helper"
)



/*  定义模型相关参数和方法  */
type OCRCardNo struct{}

func (x *OCRCardNo) Init() error {
	return nil
}

func (x *OCRCardNo) ApiPath() string {
	return "/ocr2/card_no"
}

func (x *OCRCardNo) CustomQueue() string {
	return helper.Settings.Customer["OCR_QUEUE"]
}

func (x *OCRCardNo) ApiEntry(reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Api_OCRCardNo")

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


// OCRCardNo 推理 - 不在这里实现，由 python dispatcher 实现
func (x *OCRCardNo) Infer(reqId string, reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Infer_OCRCardNo - Do nothing")

	return &map[string]interface{}{}, nil
}
