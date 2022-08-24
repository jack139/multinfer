package yhface

import (
	"fmt"
	"log"

	"multinfer/gosearch"
)

/*  定义模型相关参数和方法  */
type FaceMemory struct{}

func (x *FaceMemory) Init() error {
	return initModel()
}

func (x *FaceMemory) ApiPath() string {
	return "/private/face2/memory"
}

func (x *FaceMemory) ApiEntry(reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Api_FaceMemory")

	// 构建请求参数
	reqDataMap := map[string]interface{}{}

	return &reqDataMap, nil
}


// 推理
func (x *FaceMemory) Infer(requestId string, reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Infer_FaceMemory")

	groupId := (*reqData)["group_id"].(string)
	label := (*reqData)["label"].(string)
	action := (*reqData)["action"].(string)
	data := (*reqData)["data"].([]float32)

	log.Println(data)
	log.Println(action, groupId, label)

	switch action {
	case "add":
		gosearch.Add(groupId, label, data)

	case "remove":
		gosearch.Remove(groupId, label)

	default:
		return &map[string]interface{}{"code":9001}, fmt.Errorf("Unknown action")
	}

	return &map[string]interface{}{"msg":"ok"}, nil
}
