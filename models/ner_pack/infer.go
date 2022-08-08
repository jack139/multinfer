package ner_pack

import (
	"log"
	//"strings"
	"unicode/utf8"

	"github.com/buckhx/gobert/tokenize"
	"github.com/buckhx/gobert/tokenize/vocab"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	//"github.com/aclements/go-gg/generic/slice"

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
	//warmup()

	return nil
}

/* 判断是否是英文字符 */
func isAlpha(c byte) bool {
	return (c>=65 && c<=90) || (c>=97 && c<=122)
}

// 模型推理
func modleInfer(text string) ([]nerStruct, int, error){
	seqLen := MaxSeqLength
	if utf8.RuneCountInString(text) < MaxSeqLength {
		seqLen = utf8.RuneCountInString(text) + 2
	}

	//log.Println("seqlen: ", seqLen)

	tkz := tokenize.NewTokenizer(voc)
	ff := tokenize.FeatureFactory{Tokenizer: tkz, SeqLen: int32(seqLen)}
	// 拼接输入 
	input_tokens := text
	// 获取 token 向量,  "[CLS]" + text + "[SEP]"
	f := ff.Feature(input_tokens)

	log.Println(input_tokens)
	log.Println(f.TokenIDs)

	new_tids := make([]float32, len(f.TokenIDs))
	for i, v := range f.TokenIDs {
		new_tids[i] = float32(v)
	}
	tids, err := tf.NewTensor([][]float32{new_tids})
	if err != nil {
		return nil, 9002, err
	}
	new_sids := make([]float32, len(f.TypeIDs))
	for i, v := range f.TypeIDs {
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
	log.Println("Shape", res[0].Shape())

	scores := res[0].Value().([][][][]float32)[0]

	var result []nerStruct

	for l:=0;l<6;l++ {
		for start:=0;start<seqLen;start++ {
			for end:=0;end<seqLen;end++ {
				if scores[l][start][end] > 0 {
					log.Println(l, start, end, f.Tokens[start], f.Tokens[end])

					var v string
					for i:=start;i<end+1;i++ { v += f.Tokens[i] }
					result = append(result, nerStruct{start, labelName[l], v})
				}
			}
		}
	}

	return result, 0, nil

}

func warmup(){
	r, _, err := modleInfer(
		"什么是深度学习？",
	)
	if err==nil {
		log.Printf("warmup: %s", r)
	} else {
		log.Printf("warmup fail: %s", err.Error())
	}
}