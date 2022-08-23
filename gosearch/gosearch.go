package gosearch

import (
	"log"
	"errors"
	"strconv"

	"github.com/jack139/go-infer/helper"

	"multinfer/gosearch/facelib"
)

var (
	ThreshHold float32
)

func LoadFaceData() error {
	// DB 链接 参数
	facelib.MongoURL = helper.Settings.Customer["FACE_mongo_uri"]
	facelib.MongoUser = helper.Settings.Customer["FACE_mongo_user"]
	facelib.MongoPwd = helper.Settings.Customer["FACE_mongo_pwd"]

	// 链接数据库
	if !facelib.ConnectDb() {
		panic(errors.New("Connect DB fail!"))
	}

	// 读取特征数据
	err := facelib.ReadData(helper.Settings.Customer["FACE_GroupIdList"])
	if err != nil {
		return err
	}

	// 初始化参数
	value, err := strconv.ParseFloat(helper.Settings.Customer["FACE_DistanceThreshold"], 32)
	if err != nil {
		return err
	}
	ThreshHold = float32(value)

	facelib.LimitFace, err = strconv.Atoi(helper.Settings.Customer["FACE_LimitFace"])
	if err != nil {
		return err
	}

	facelib.GONUM, err = strconv.Atoi(helper.Settings.Customer["FACE_Gonum"])
	if err != nil {
		return err
	}

	return nil
}

// 根据特征值搜索
func Search(groupId string, testVec []float32) *map[string]interface{} {
	label, min := facelib.Search(groupId, testVec)
	log.Println("gosearch-Search: ", groupId, ThreshHold, label, min)

	if min < ThreshHold && label!="__BLANK__" { // __BLANK__ 说明特征已动态删除
		resultMap := map[string]interface{}{
			"label": label,
			"score": min,
		}
		return &resultMap
	} else {
		return nil
	}
}

// 内存中添加特征
func Add(groupId, label string, testVec []float32) {
	log.Println("gosearch-Add: ", groupId, label)
	facelib.AddNewData(groupId, label, testVec)

	return
}

// 内存中删除特征
func Remove(groupId, label string) {
	log.Println("gosearch-Remove: ", groupId, label)
	facelib.RemoveData(groupId, label)

	return
}
