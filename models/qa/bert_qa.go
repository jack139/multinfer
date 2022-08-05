package qa

import (
	"fmt"
	"log"
	"strings"

	"github.com/buckhx/gobert/tokenize"
	"github.com/buckhx/gobert/tokenize/vocab"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/aclements/go-gg/generic/slice"

	"github.com/jack139/go-infer/helper"
)

const (
	MaxSeqLength = 512
)

/* 训练好的模型权重 */
var (
	m *tf.SavedModel
	voc vocab.Dict
)

/* 初始化模型 */
func initModel() error {
	var err error
	voc, err = vocab.FromFile(helper.Settings.Customer["BertVocabPath"])
	if err != nil {
		return err
	}
	m, err = tf.LoadSavedModel(helper.Settings.Customer["BertModelPath"], []string{"train"}, nil)
	if err != nil {
		return err
	}

	return nil
}

/* 判断是否是英文字符 */
func isAlpha(c byte) bool {
	return (c>=65 && c<=90) || (c>=97 && c<=122)
}


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

	tkz := tokenize.NewTokenizer(voc)
	ff := tokenize.FeatureFactory{Tokenizer: tkz, SeqLen: MaxSeqLength}
	// 拼接输入
	input_tokens := question + tokenize.SequenceSeparator + corpus
	// 获取 token 向量
	f := ff.Feature(input_tokens)

	tids, err := tf.NewTensor([][]int32{f.TokenIDs})
	if err != nil {
		return &map[string]interface{}{"code":9002}, err
	}
	new_mask := make([]float32, len(f.Mask))
	for i, v := range f.Mask {
		new_mask[i] = float32(v)
	}
	mask, err := tf.NewTensor([][]float32{new_mask})
	if err != nil {
		return &map[string]interface{}{"code":9003}, err
	}
	sids, err := tf.NewTensor([][]int32{f.TypeIDs})
	if err != nil {
		return &map[string]interface{}{"code":9004}, err
	}

	res, err := m.Session.Run(
		map[tf.Output]*tf.Tensor{
			m.Graph.Operation("input_ids").Output(0):      tids,
			m.Graph.Operation("input_mask").Output(0):     mask,
			m.Graph.Operation("segment_ids").Output(0):    sids,
		},
		[]tf.Output{
			m.Graph.Operation("finetune_mrc/Squeeze").Output(0),
			m.Graph.Operation("finetune_mrc/Squeeze_1").Output(0),
		},
		nil,
	)
	if err != nil {
		return &map[string]interface{}{"code":9005}, err
	}

	st := slice.ArgMax(res[0].Value().([][]float32)[0])
	ed := slice.ArgMax(res[1].Value().([][]float32)[0])
	//fmt.Println(st, ed)
	if ed<st{ // ed 小于 st 说明未找到答案
		st = 0
		ed = 0
	}
	//ans = strings.Join(f.Tokens[st:ed+1], "")

	// 处理token中的英文，例如： 'di', '##st', '##ri', '##bu', '##ted', 're', '##pr', '##ese', '##nt', '##ation',
	var ans string
	for i:=st;i<ed+1;i++ {
		if len(f.Tokens[i])>0 && isAlpha(f.Tokens[i][0]){ // 英文开头，加空格
			ans += " "+f.Tokens[i]
		} else if strings.HasPrefix(f.Tokens[i], "##"){ // ##开头，是英文中段，去掉##
			ans += f.Tokens[i][2:]
		} else {
			ans += f.Tokens[i]
		}
	}

	if strings.HasPrefix(ans, "[CLS]") || strings.HasPrefix(ans, "[SEP]") {
		return &map[string]interface{}{"answer":""}, nil
	} else {
		return &map[string]interface{}{"answer":ans}, nil // 找到答案
	}
}
