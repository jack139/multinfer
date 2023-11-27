package talk2ui

import (
	"fmt"
	"log"
	"time"
	"crypto/md5"
	"math/rand"

	"github.com/jack139/go-infer/helper"
	"github.com/jack139/go-infer/types"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

/* 产生随机串 */
func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

/* 产生 request id */
func generateRequestId() string {
	//year, month, day := time.Now().Date()
	dateStr := time.Now().Format("20060102150405")
	h := md5.New()
	h.Write([]byte(randSeq(10)))
	sum := h.Sum(nil)
	md5Str := fmt.Sprintf("wav%s%x", dateStr, sum)
	return md5Str
}


/*  定义模型相关参数和方法  */
type Wav2Order struct{}

func (x *Wav2Order) Init() error {
	return nil
}

func (x *Wav2Order) ApiPath() string {
	return "/talk2ui/wav2order"
}

func (x *Wav2Order) CustomQueue() string {
	return helper.Settings.Customer["BERT_QUEUE"] // 使用 Bert 队列进行最后的识别
}

func (x *Wav2Order) ApiEntry(reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Api_Wav2Order")

	// 检查参数
	wavData, ok := (*reqData)["wav_data"].(string)
	if !ok {
		return &map[string]interface{}{"code":9101}, fmt.Errorf("need wav_data")
	}

	// 先调用 Wav2Text() 获取文本

	// 构建请求参数
	reqDataMap := map[string]interface{}{
		"text": "",
	}

	// 构建 wav2text 的 请求参数
	reqDataMap2 := map[string]interface{}{
		"wav_data": wavData,
	}

	// 手工调用 wav2text
	for mIndex := range types.ModelList {
		if types.ModelList[mIndex].ApiPath() == "/talk2ui/wav2text" {
			// 临时的请求 id
			requestId := generateRequestId()

			// 构建队列请求参数
			reqQueueDataMap := map[string]interface{}{
				"api": types.ModelList[mIndex].ApiPath(),
				"params": reqDataMap2,
			}


			// 注册消息队列，在发redis消息前注册, 防止消息漏掉
			pubsub := helper.Redis_subscribe(requestId)
			defer pubsub.Close()

			// 发 请求消息
			queueName := types.ModelList[mIndex].CustomQueue()
			err := helper.Redis_publish_request(requestId, queueName, &reqQueueDataMap)
			if err!=nil {
				return &map[string]interface{}{"code":9102}, fmt.Errorf(
					helper.Settings.ErrCode.SENDMSG_FAIL["msg"].(string) + " : " + err.Error())
			}

			// 收 结果消息
			respData, err := helper.Redis_sub_receive(pubsub)
			if err!=nil {
				return &map[string]interface{}{"code":9103}, fmt.Errorf(
					helper.Settings.ErrCode.RECVMSG_FAIL["msg"].(string) + " : " + err.Error())
			}

			// code==0 提交成功
			if (*respData)["code"].(float64)!=0 { 
				return &map[string]interface{}{"code":9104}, fmt.Errorf(
					(*respData)["msg"].(string))
			}

			result := ((*respData)["data"].(map[string]interface{}))["result"].(map[string]interface{})
			reqDataMap["text"] = result["text"].(string)

		}
	}

	//log.Println(reqDataMap)

	return &reqDataMap, nil
}


// Wav2Order 推理 - 不在这里实现，由 python dispatcher 实现
func (x *Wav2Order) Infer(reqId string, reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Infer_Wav2Order - Do nothing")

	return &map[string]interface{}{}, nil
}
