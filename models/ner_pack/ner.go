package ner_pack

import (
	"fmt"
	"log"
	"sort"
)



/*  定义模型相关参数和方法  */
type NER struct{}

/* 用于分割文本 */
var seperators = []string{"；", "，", "。", ",", "）", "、", ";"}


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


func findSeperator(a string) bool {
	for _, b := range seperators {
		if b == a {
			return true
		}
	}
	return false
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
		if len(text[posOffset:]) > MaxSeqLength {
			var n int
			for n=MaxSeqLength;n>0;n-- {
				if findSeperator(string(text[posOffset:][n])) {
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

		log.Println(ans)

		if ans!=nil {
			result = append(result, ans...)
		}

		posOffset += len(text1)
	}

	sort.Slice(result, func(i, j int) bool { return result[i].startPos < result[j].startPos })

	log.Println(result)

	var result2 []map[string]interface{}
	if result!=nil {
		for i := range result{
			result2 = append(result2, map[string]interface{}{
				"start_pos" : result[i].startPos,
				"label" : result[i].label,
				"value" : result[i].value,
			})
		}
	}

	return &map[string]interface{}{"entities":result2}, nil
}
