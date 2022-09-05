package yhface

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"encoding/base64"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/jack139/go-infer/helper"

	"multinfer/models/yhface/gosearch"
	"multinfer/models/yhface/gosearch/facelib"
)

/*  定义模型相关参数和方法  */
type FaceFeatures struct{}

func (x *FaceFeatures) Init() error {
	return initModel()
}

func (x *FaceFeatures) ApiPath() string {
	return "__noapi__/features"
}

func (x *FaceFeatures) ApiEntry(reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Api_FaceFeatures")

	// 构建请求参数
	reqDataMap := map[string]interface{}{}

	return &reqDataMap, nil
}


// 推理
func (x *FaceFeatures) Infer(requestId string, reqData *map[string]interface{}) (*map[string]interface{}, error) {
	log.Println("Infer_FaceFeatures")

	imageBase64 := (*reqData)["image"].(string)
	groupId := (*reqData)["group_id"].(string)
	userId := (*reqData)["user_id"].(string)
	faceId := (*reqData)["face_id"].(string)

	// 解码base64
	image, err  := base64.StdEncoding.DecodeString(imageBase64)
	if err!=nil {
		return &map[string]interface{}{"code":9901}, err
	}

	// 检查图片大小
	maxSize, _ := strconv.Atoi(helper.Settings.Customer["FACE_MAX_IMAGE_SIZE"])
	if len(image) > maxSize {
		return &map[string]interface{}{"code":9002}, fmt.Errorf("图片数据太大")
	}

	// 模型推理
	feat, _, normFace, code, err := featuresInfer(image)
	if err != nil {
		return &map[string]interface{}{"code":code}, err
	}

	if feat==nil {  // 未检测到人脸
		return &map[string]interface{}{"features": feat}, nil
	}

	// 正则化
	feat, err = norm(feat)
	if err != nil {
		return &map[string]interface{}{"code":9005}, err
	}

	// 保存请求图片和结果
	saveBackLog(requestId, normFace, []byte(fmt.Sprintf("%v %v", userId, faceId)))


	// 更新到内存
	gosearch.Add(groupId, userId, feat)

	// 更新到DB

	// 检测数据库连接
	if !facelib.Ping() {
		return &map[string]interface{}{"code":9008}, fmt.Errorf("DB connection problem.")
	}

	// normFace 转换 为 字节流
	cropByte, err := image2bytes(normFace)
	if err != nil {
		return &map[string]interface{}{"code":9006}, err
	}

	// 更新人脸信息
	database := facelib.Client.Database("face_db")
	collFace := database.Collection("faces")

	opts := options.Update().SetUpsert(false)
	objID, _ := primitive.ObjectIDFromHex(faceId)
	filter := bson.D{{"_id", objID}}
	update := bson.M{ "$set": bson.M{ 
			"encodings": bson.M{ 
				"arc": bson.M{ 
					"None": feat,
				},
			},
			"image": cropByte,
		},
	}

	_, err = collFace.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return &map[string]interface{}{"code":9009}, err
	}

	// 未返回全部特征值，只为区别于 nil
	return &map[string]interface{}{"features":feat[:10]}, nil 
}
