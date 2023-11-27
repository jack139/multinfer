package keras_qa

import (
	"fmt"
	"log"
)



/*  定义模型相关参数和方法  */
type AlbertQA struct{}

func (x *AlbertQA) Init() error {
	return initModel()
}

func (x *AlbertQA) ApiPath() string {
	return "/api/albert_qa"
}

func (x *AlbertQA) CustomQueue() string {
	return ""
}

func (x *AlbertQA) ApiEntry(reqData *map[string]interface{}) (*map[string]interface{}, error) {
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
func (x *AlbertQA) Infer(reqId string, reqData *map[string]interface{}) (*map[string]interface{}, error) {
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
