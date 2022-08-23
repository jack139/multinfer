package facelib

import (
	//"log"
	//"math"
)

// 管道返回值
type ChResult struct {
    val float32
    label uint32
}

var (
	GONUM = int(8) // go routine 数量, 建议与cpu核数一致
)

// 计算欧式距离, 不开根号
func edist(x []float32, y []float32) float32 {
	var sum float32
	for i:=0;i<len(x);i++ { // plus
		sum += (x[i]-y[i])*(x[i]-y[i])
	}
	//sum = float32(math.Sqrt(float64(sum)))

	return sum
}



// 计算余弦相似性
func cosine(x []float32, y []float32) float32 {
	var sum float32
	for i:=0;i<len(x);i++ { // plus
		sum += (x[i] * y[i])
	}

	return -sum
}


// 用X中最后1个向量做测试
func findMin(group string, target []float32, start int, end int, ch chan ChResult) {
	//log.Println("-->", start, end)
	var min1 float32
	var label uint32
	min1 = 999999999.0
	for i:=start; i<end; i++ {
		if labelName[group][y[group][i]-1]=="__BLANK__" {
			continue // 跳过已删除的
		}
		dist := cosine(X[group][i], target)
		//log.Printf("%.8f ", dist)
		if dist<min1 {
			min1 = dist
			label = y[group][i]
		}
	}
	//log.Println()

	res := new(ChResult)
	res.val = min1
	res.label = label
	ch<- *res
}


// 使用特征向量进行检索，返回最近的 label name 和距离值 
func Search(group string, target []float32) (string, float32) {
	var min float32
	var label uint32
	var seg int
	min = 99999.0
	label = 0
	channel := make([]chan ChResult, GONUM)

	N := len(X[group])
	seg = N/GONUM

	for i:=0;i<GONUM;i++ {
		channel[i] = make(chan ChResult)
		var end int
		if i+1==GONUM { 
			end = N // 对seg计算整除有余数的情况
		} else {
			end = i*seg+seg
		}
		go findMin(group, target, i*seg, end, channel[i])
	}

	// 取得返回结果
	for _, rc := range(channel) {
		t := <-rc
		if t.val<min {
			min = t.val
			label = t.label
		}
	}

	return labelName[group][label-1], min

}


