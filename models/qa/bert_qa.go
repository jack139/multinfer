package qa

import (
	"fmt"
	"log"
)



/*  定义模型相关参数和方法  */
type BertQA struct{}

func (x *BertQA) Init() error {
	return initModel()
}

func (x *BertQA) ApiPath() string {
	return "/api/bert_qa"
}

func (x *BertQA) ApiEntry(reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Api_BertQA")

	// 检查参数
	corpus, ok := (*reqData)["corpus"].(string)
	if !ok {
		return &map[string]interface{}{"code":9101}, fmt.Errorf("need corpus")
	}

	question, ok := (*reqData)["question"].(string)
	if !ok {
		return &map[string]interface{}{"code":9102}, fmt.Errorf("need question")
	}

	// 构建请求参数
	reqDataMap := map[string]interface{}{
		"corpus": corpus,
		"question": question,
	}

	return &reqDataMap, nil
}


// Bert 推理
func (x *BertQA) Infer(reqId string, reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Infer_BertQA")

	corpus := (*reqData)["corpus"].(string)
	question := (*reqData)["question"].(string)
	//log.Printf("Corpus: %s\tQuestion: %s", corpus, question)

	ans, code, err := modleInfer(corpus, question)
	if err != nil {
		return &map[string]interface{}{"code":code}, err
	}

	return &map[string]interface{}{"answer":ans}, nil
}
