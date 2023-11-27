package talk2ui

import (
	"fmt"
	"log"
)



/*  定义模型相关参数和方法  */
type Wav2Order struct{}

func (x *Wav2Order) Init() error {
	return initModel()
}

func (x *Wav2Order) ApiPath() string {
	return "/tail2ui/text2order"
}

func (x *Wav2Order) CustomQueue() string {
	return ""
}

func (x *Wav2Order) ApiEntry(reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Api_Wav2Order")

	// 检查参数
	wavData, ok := (*reqData)["wav_data"].(string)
	if !ok {
		return &map[string]interface{}{"code":9101}, fmt.Errorf("need wav_data")
	}

	// 构建请求参数
	reqDataMap := map[string]interface{}{
		"wav_data": wavData,
	}

	return &reqDataMap, nil
}


// Wav2Order 推理 - 不在这里实现，由 python dispatcher 实现
func (x *Wav2Order) Infer(reqId string, reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Infer_Wav2Order - Do nothing")

	return &map[string]interface{}{}, nil
}
