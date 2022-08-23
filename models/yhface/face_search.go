package yhface

import (
	"fmt"
	"log"
	"strconv"
	"encoding/base64"

	"github.com/jack139/go-infer/helper"

	"multinfer/gosearch"
)

/*  定义模型相关参数和方法  */
type FaceSearch struct{}

func (x *FaceSearch) Init() error {
	return initModel()
}

func (x *FaceSearch) ApiPath() string {
	return "/face2/search"
}

func (x *FaceSearch) ApiEntry(reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Api_FaceSearch")

	// 检查参数
	imageBase64, ok := (*reqData)["image"].(string)
	if !ok {
		return &map[string]interface{}{"code":9001}, fmt.Errorf("need image")
	}

	var groupId, userId, mobileTail string
	var maxUser float64

	groupId, ok = (*reqData)["group_id"].(string)
	if !ok {
		groupId = "DEFAULT"
	}

	userId, ok = (*reqData)["user_id"].(string)
	if !ok {
		userId = ""
	}

	mobileTail, ok = (*reqData)["mobile_tail"].(string)
	if !ok {
		mobileTail = ""
	}

	maxUser, ok = (*reqData)["max_user_num"].(float64)
	if !ok {
		maxUser = 5
	}

	// 构建请求参数
	reqDataMap := map[string]interface{}{
		"image": imageBase64,
		"group_id" : groupId,
		"user_id" : userId,
		"mobile_tail" : mobileTail,
		"max_user_num" : maxUser,
	}

	return &reqDataMap, nil
}


// 推理
func (x *FaceSearch) Infer(requestId string, reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Infer_FaceSearch")

	imageBase64 := (*reqData)["image"].(string)
	groupId := (*reqData)["group_id"].(string)
	//userId := (*reqData)["user_id"].(string)
	//mobileTail := (*reqData)["mobile_tail"].(string)
	//maxUser := (*reqData)["max_user_num"].(float64)

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
	feat, code, err := featuresInfer(image)
	if err != nil {
		return &map[string]interface{}{"code":code}, err
	}

	// 正则化
	feat, err = norm(feat)
	if err != nil {
		return &map[string]interface{}{"code":9005}, err
	}

	r := gosearch.Search(groupId, feat)

	// 保存请求图片和结果
	saveBackLog(requestId, image, []byte(fmt.Sprintf("%v", r)))

	return &map[string]interface{}{"result":r}, nil
}
