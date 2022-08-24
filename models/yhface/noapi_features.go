package yhface

import (
	"fmt"
	"log"
	"strconv"
	"encoding/base64"

	"github.com/jack139/go-infer/helper"
)

/*  定义模型相关参数和方法  */
type FaceFeatures struct{}

func (x *FaceFeatures) Init() error {
	return initModel()
}

func (x *FaceFeatures) ApiPath() string {
	return "/private/face2/features"
}

func (x *FaceFeatures) ApiEntry(reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Api_FaceFeatures")

	// 构建请求参数
	reqDataMap := map[string]interface{}{}

	return &reqDataMap, nil
}


// 推理
func (x *FaceFeatures) Infer(requestId string, reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Infer_FaceFeatures")

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
	feat, _, code, err := featuresInfer(image)
	if err != nil {
		return &map[string]interface{}{"code":code}, err
	}

	if feat==nil {  // 未检测到人脸
		return &map[string]interface{}{"features": feat}, nil
	}

	// 正则化
	feat, err = norm(feat)
	if err != nil {
		return &map[string]interface{}{"code":9005}, err
	}

	// 保存请求图片和结果
	saveBackLog(requestId, image, []byte(fmt.Sprintf("%v", feat[:10])))

	return &map[string]interface{}{"features":feat}, nil
}
