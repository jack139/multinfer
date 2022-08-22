package yhface

import (
	//"os"
	//"time"
	"fmt"
	"log"
	"strconv"
	"encoding/base64"
	//"io/ioutil"

	"github.com/jack139/go-infer/helper"
)

/*  定义模型相关参数和方法  */
type FaceCheck struct{}

func (x *FaceCheck) Init() error {
	return initModel()
}

func (x *FaceCheck) ApiPath() string {
	return "/face2/check"
}

func (x *FaceCheck) ApiEntry(reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Api_FaceCheck")

	// 检查参数
	imageBase64, ok := (*reqData)["image"].(string)
	if !ok {
		return &map[string]interface{}{"code":9001}, fmt.Errorf("need image")
	}

	// 构建请求参数
	reqDataMap := map[string]interface{}{
		"image": imageBase64,
	}

	return &reqDataMap, nil
}


// 推理
func (x *FaceCheck) Infer(requestId string, reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Infer_FaceCheck")

	imageBase64 := (*reqData)["image"].(string)

	// 解码base64
	image, err  := base64.StdEncoding.DecodeString(imageBase64)
	if err!=nil {
		return &map[string]interface{}{"code":9901}, err
	}

	// 检查图片大小
	maxSize, _ := strconv.Atoi(helper.Settings.Customer["FACE_MAX_IMAGE_SIZE"])
	if len(image) > maxSize {
		return &map[string]interface{}{"code":9002}, fmt.Errorf("图片数据太大")
	}

	// 模型推理
	r, code, err := locateInfer(image)
	if err != nil {
		return &map[string]interface{}{"code":code}, err
	}

	log.Println("face num--> ", len(r))

	// 保存请求图片和识别结果（文件名中体现结果）
	/*
	if helper.Settings.Customer["FACE_SAVE_IMAGE"] == "1" {
		output_dir := fmt.Sprintf("%s/%s", 
			helper.Settings.Customer["FACE_SAVE_IMAGE_PATH"], 
			time.Now().Format("20060102"))
		err = os.Mkdir(output_dir, 0755) // 建日志目录， 日期 做子目录
		if err == nil || os.IsExist(err) { // 不处理错误
			_ = ioutil.WriteFile(fmt.Sprintf("%s/%s_%s.jpg", output_dir, requestId, r), image, 0644)
		} else {
			log.Println("ERROR when saving log: ", err.Error())
		}
	}
	*/

	return &map[string]interface{}{"has_face": len(r)>0}, nil
}
