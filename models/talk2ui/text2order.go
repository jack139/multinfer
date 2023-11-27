package talk2ui

import (
	"fmt"
	"log"
)



/*  定义模型相关参数和方法  */
type Text2Order struct{}

func (x *Text2Order) Init() error {
	return initModel()
}

func (x *Text2Order) ApiPath() string {
	return "/tail2ui/text2order"
}

func (x *Text2Order) CustomQueue() string {
	return ""
}

func (x *Text2Order) ApiEntry(reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Api_Text2Order")

	// 检查参数
	text, ok := (*reqData)["text"].(string)
	if !ok {
		return &map[string]interface{}{"code":9101}, fmt.Errorf("need text")
	}

	// 构建请求参数
	reqDataMap := map[string]interface{}{
		"text": text,
	}

	return &reqDataMap, nil
}


// Text2Order 推理 - 不在这里实现，由 python dispatcher 实现
func (x *Text2Order) Infer(reqId string, reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Infer_Text2Order - Do nothing")

	return &map[string]interface{}{}, nil
}
