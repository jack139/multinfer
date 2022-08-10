package ner_pack

import (
	"log"
	"strings"

	"github.com/buckhx/gobert/tokenize"
	"github.com/buckhx/gobert/tokenize/vocab"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"

	"github.com/jack139/go-infer/helper"
)

/* 结果返回结构 */
type nerStruct struct{
	startPos int
	label, value string
}

const (
	MaxSeqLength = 512
)

/* 训练好的模型权重 */
var (
	m *tf.SavedModel
	voc vocab.Dict

	labelName = []string{"检验和检查", "治疗和手术", "疾病和诊断", "症状和体征", "药物", "解剖部位"}

	/* 用于分割文本 */
	seperators = []string{"；", "，", "。", ",", "、", ";", "）", ")"}
)

/* 初始化模型 */
func initModel() error {
	var err error
	voc, err = vocab.FromFile(helper.Settings.Customer["NerPackVocabPath"])
	if err != nil {
		return err
	}
	m, err = tf.LoadSavedModel(helper.Settings.Customer["NerPackModelPath"], []string{"train"}, nil)
	if err != nil {
		return err
	}

	// 模型热身
	warmup()

	return nil
}

// 找字符串
func strInStrings(a string, s []string) bool {
	for _, b := range s {
		if b == a {
			return true
		}
	}
	return false
}

/* 判断是否是英文字符 */
func isAlpha(c byte) bool {
	return (c>=65 && c<=90) || (c>=97 && c<=122)
}

// 模型推理
func modleInfer(text string, posOffset int) ([]nerStruct, int, error){
	tkz := tokenize.NewTokenizer(voc)
	ff := tokenize.FeatureFactory{Tokenizer: tkz, SeqLen: MaxSeqLength}
	// 拼接输入 
	input_tokens := text
	// 获取 token 向量,  "[CLS]" + text + "[SEP]"
	f := ff.Feature(input_tokens)

	//log.Println(input_tokens, len([]rune(input_tokens)))
	//log.Println(f.TokenIDs, len(f.TokenIDs))
	//log.Println(f.Tokens, len(f.Tokens))
	//log.Println(f.Count())

	// 还原token与原始文本的对照关系
	orig_pos := 0
	orig_text := []rune(input_tokens)
	orig_token := make([]nerStruct, f.Count())
	for i, v := range f.Tokens[:f.Count()] {
		if v=="[UNK]" { // 位置token, 没有在vocab里的字
			orig_token[i] = nerStruct{orig_pos, "", string(orig_text[orig_pos])}
			orig_pos++
			continue
		}
		if strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]") { 
			// 其他bert标签 [CLS] [SEP] [PAD] [unused] 等
			continue 
		} 
		if strings.HasPrefix(v, "##") { // 英文, ##前不用处理空格的问题
			orig_token[i] = nerStruct{orig_pos, "", string(orig_text[orig_pos:orig_pos+len(v)-2])}
			orig_pos = orig_pos + len(v) - 2
		} else {
			// 匹配字符串，可以跳过可能的空格
			for strings.ToLower(
					string(orig_text[orig_pos:orig_pos+len([]rune(v))]),
				) != strings.ToLower(v) { orig_pos++ }
			orig_token[i] = nerStruct{orig_pos, "", string(orig_text[orig_pos:orig_pos+len([]rune(v))])}
			orig_pos = orig_pos + len([]rune(v))
		}
	}

	//log.Println(orig_token, len(orig_token))


	new_tids := make([]float32, f.Count())
	for i, v := range f.TokenIDs[:f.Count()] {
		new_tids[i] = float32(v)
	}
	tids, err := tf.NewTensor([][]float32{new_tids})
	if err != nil {
		return nil, 9002, err
	}
	new_sids := make([]float32, f.Count())
	for i, v := range f.TypeIDs[:f.Count()] {
		new_sids[i] = float32(v)
	}
	sids, err := tf.NewTensor([][]float32{new_sids})
	if err != nil {
		return nil, 9004, err
	}

	res, err := m.Session.Run(
		map[tf.Output]*tf.Tensor{
			m.Graph.Operation("Input-Token").Output(0):      tids,
			m.Graph.Operation("Input-Segment").Output(0):    sids,
		},
		[]tf.Output{
			m.Graph.Operation("efficient_global_pointer_1/sub_3").Output(0),
		},
		nil,
	)
	if err != nil {
		return nil, 9005, err
	}

	//log.Printf("%v", res[0].Value().([][][][]float32))
	log.Println("Result shape: ", res[0].Shape())

	scores := res[0].Value().([][][][]float32)[0]

	var result []nerStruct

	for l:=0;l<6;l++ {
		for start:=0;start<f.Count();start++ {
			for end:=0;end<f.Count();end++ {
				if scores[l][start][end] > 0 {
					//log.Println(l, start, end, f.Tokens[start], f.Tokens[end])

					// 处理token中的英文，例如： 'di', '##st', '##ri', '##bu', '##ted', 're', '##pr', '##ese', '##nt', '##ation',
					var ans string
					for i:=start;i<end+1;i++ {
						ans += orig_token[i].value
					}
					// start 是相对位置，所以要加上 posOffset, 开头有[CLS],所以减1
					result = append(result, nerStruct{orig_token[start].startPos+posOffset, labelName[l], ans})
				}
			}
		}
	}

	return result, 0, nil

}

func warmup(){
	r, _, err := modleInfer(
		"于当地行胃镜检查并行病理检查示", 0,
	)
	if err==nil {
		log.Printf("warmup: %s", r)
	} else {
		log.Printf("warmup fail: %s", err.Error())
	}
}