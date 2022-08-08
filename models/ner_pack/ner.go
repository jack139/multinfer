package ner_pack

import (
	"fmt"
	"log"
)



/*  定义模型相关参数和方法  */
type NER struct{}

func (x *NER) Init() error {
	return initModel()
}

func (x *NER) ApiPath() string {
	return "/ner/ner"
}

func (x *NER) ApiEntry(reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Api_NER")

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


// NER 推理
func (x *NER) Infer(reqId string, reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Infer_NER")

	text := (*reqData)["text"].(string)

	ans, code, err := modleInfer(text)
	if err != nil {
		return &map[string]interface{}{"code":code}, err
	}

	log.Println(ans)

	return &map[string]interface{}{"answer":ans}, nil
}
