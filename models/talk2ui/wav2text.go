package talk2ui

import (
	"fmt"
	"log"
	"strconv"

	"github.com/jack139/go-infer/helper"
)



/*  定义模型相关参数和方法  */
type Wav2Text struct{}

func (x *Wav2Text) Init() error {
	return nil
}

func (x *Wav2Text) ApiPath() string {
	return "/talk2ui/wav2text"
}

func (x *Wav2Text) CustomQueue() string {
	return helper.Settings.Customer["ASR_QUEUE"]
}

func (x *Wav2Text) ApiEntry(reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Api_Wav2Text")

	// 检查参数
	wavData, ok := (*reqData)["wav_data"].(string)
	if !ok {
		return &map[string]interface{}{"code":9101}, fmt.Errorf("need wav_data")
	}

	// 检查数据大小
	maxSize, _ := strconv.Atoi(helper.Settings.Customer["WAV_MAX_IMAGE_SIZE"])
	if len(wavData) > maxSize {
		return &map[string]interface{}{"code":9002}, fmt.Errorf("语音数据太大")
	}

	// 构建请求参数
	reqDataMap := map[string]interface{}{
		"wav_data": wavData,
	}

	return &reqDataMap, nil
}


// Wav2Text 推理 - 不在这里实现，由 python dispatcher 实现
func (x *Wav2Text) Infer(reqId string, reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Infer_Wav2Text - Do nothing")

	return &map[string]interface{}{}, nil
}
