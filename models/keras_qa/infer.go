package keras_qa

import (
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
	voc, err = vocab.FromFile(helper.Settings.Customer["ALBertVocabPath"])
	if err != nil {
		return err
	}
	m, err = tf.LoadSavedModel(helper.Settings.Customer["ALBertModelPath"], []string{"train"}, nil)
	if err != nil {
		return err
	}

	// 模型热身
	warmup()

	return nil
}

/* 判断是否是英文字符 */
func isAlpha(c byte) bool {
	return (c>=65 && c<=90) || (c>=97 && c<=122)
}

// 模型推理
func modleInfer(corpus, question string) (string, int, error){

	tkz := tokenize.NewTokenizer(voc)
	ff := tokenize.FeatureFactory{Tokenizer: tkz, SeqLen: MaxSeqLength}
	// 拼接输入
	input_tokens := question + tokenize.SequenceSeparator + corpus
	// 获取 token 向量
	f := ff.Feature(input_tokens)

	new_tids := make([]float32, len(f.TokenIDs))
	for i, v := range f.TokenIDs {
		new_tids[i] = float32(v)
	}
	tids, err := tf.NewTensor([][]float32{new_tids})
	if err != nil {
		return "", 9002, err
	}
	//new_mask := make([]float32, len(f.Mask))
	//for i, v := range f.Mask {
	//	new_mask[i] = float32(v)
	//}
	//mask, err := tf.NewTensor([][]float32{new_mask})
	//if err != nil {
	//	return "", 9003, err
	//}
	new_sids := make([]float32, len(f.TypeIDs))
	for i, v := range f.TypeIDs {
		new_sids[i] = float32(v)
	}
	sids, err := tf.NewTensor([][]float32{new_sids})
	if err != nil {
		return "", 9004, err
	}

	res, err := m.Session.Run(
		map[tf.Output]*tf.Tensor{
			m.Graph.Operation("Input-Token").Output(0):      tids,
			//m.Graph.Operation("input_mask").Output(0):     mask,
			m.Graph.Operation("Input-Segment").Output(0):    sids,
		},
		[]tf.Output{
			m.Graph.Operation("permute/transpose").Output(0),
			//m.Graph.Operation("finetune_mrc/Squeeze_1").Output(0),
		},
		nil,
	)
	if err != nil {
		return "", 9005, err
	}

	st := slice.ArgMax(res[0].Value().([][][]float32)[0][0])
	ed := slice.ArgMax(res[0].Value().([][][]float32)[0][1])
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
		return "", 0, nil
	} else {
		return ans, 0, nil // 找到答案
	}

}

func warmup(){
	r, _, err := modleInfer(
		"深度学习（英语：deep learning）是机器学习的分支，是一种以人工神经网络为架构，对资料进行表征学习的算法。",
		"什么是深度学习？",
	)
	if err==nil {
		log.Printf("warmup: %s", r)
	} else {
		log.Printf("warmup fail: %s", err.Error())
	}
}