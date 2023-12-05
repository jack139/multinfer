package ner_pack

import (
	"fmt"
	"log"
	"sort"

	"github.com/jack139/go-infer/types"
)



/*  定义模型相关参数和方法  */
type NER struct{ types.Base }


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

	original_text := (*reqData)["text"].(string)
	text := []rune(original_text)

	var result []nerStruct
	var text1 []rune

	posOffset := 0
	for posOffset < len(text) {
		if len(text[posOffset:]) > MaxSeqLength { // 文本长度大于最大限制，则分段处理
			var n int
			for n=MaxSeqLength;n>0;n-- {
				if strInStrings(string(text[posOffset:][n]), seperators) {
					break
				}
			}
			text1 = text[posOffset:][:n]
		} else {
			text1 = text[posOffset:]
		}

		ans, code, err := modleInfer(string(text1), posOffset)
		if err != nil {
			return &map[string]interface{}{"code":code}, err
		}

		//log.Println(ans)

		if ans!=nil {
			result = append(result, ans...)
		}

		posOffset += len(text1)
	}

	// 按起始位置排序
	sort.Slice(result, func(i, j int) bool { return result[i].startPos < result[j].startPos })

	//log.Println(result)

	// 准备返回结果
	var result2 []map[string]interface{}
	if result!=nil {
		for _, v := range result{
			result2 = append(result2, map[string]interface{}{
				"start_pos" : v.startPos,
				"label" :     v.label,
				"value" :     v.value,
			})
		}
	}

	return &map[string]interface{}{"entities":result2}, nil
}
