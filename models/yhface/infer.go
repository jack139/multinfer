package yhface

import (
	"context"
	"log"
	"fmt"
	"bytes"
	"image"
	"io/ioutil"

	"github.com/disintegration/imaging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/jack139/go-infer/helper"
	"github.com/jack139/arcface-go/arcface"

	"multinfer/models/yhface/gosearch"
	"multinfer/models/yhface/gosearch/facelib"
	"multinfer/models/yhface/fas2"
)

const vecLen = 512 // 特征向量长度

/* 训练好的模型权重 */
var (
	initOK = bool(false)
)

/* 初始化模型 */
func initModel() error {
	var err error

	if !initOK { // 模型只装入一次
		if err = arcface.LoadOnnxModel(helper.Settings.Customer["ArcfaceModelPath"]); err!=nil {
			return err
		}
		log.Println("Arcface onnx model loaded from: ", helper.Settings.Customer["ArcfaceModelPath"])

		if err = fas2.LoadOnnxModel(helper.Settings.Customer["Fas2ModelPath"]); err!=nil {
			return err
		}
		log.Println("FAS onnx model loaded from: ", helper.Settings.Customer["Fas2ModelPath"])


		// 人脸库装入内存
		if err = gosearch.LoadFaceData(); err!=nil {
			return err
		}

		// 初始化标记
		initOK = true

		// 模型热身
		warmup(helper.Settings.Customer["FACE_WARM_UP_IMAGES"])
	}

	return nil
}


func locateInfer(imageByte []byte) ([][]float32, int, error){

	// 转换为 image.Image
	reader := bytes.NewReader(imageByte)

	img, err := imaging.Decode(reader)
	if err!=nil {
		return nil, 9201,err
	}

	// 检测人脸
	dets, _, err := arcface.FaceDetect(img)
	if err != nil {
		return nil, 9202, err
	}

	return dets, 0, nil
}


func featuresInfer(imageByte []byte) ([]float32, []float32, image.Image, int, error){

	// 转换为 image.Image
	reader := bytes.NewReader(imageByte)

	img, err := imaging.Decode(reader)
	if err!=nil {
		return nil, nil, nil, 9201, err
	}

	// 检测人脸
	dets, kpss, err := arcface.FaceDetect(img)
	if err != nil {
		return nil, nil, nil, 9202, err
	}

	if len(dets)==0 {
		log.Println("No face detected.")
		return nil, nil, nil, 0, nil
	}

	// 只返回第一个人脸的特征
	features, normFace, err := arcface.FaceFeatures(img, kpss[0])
	if err != nil {
		return nil, nil, nil, 9203, err
	}

	return features, dets[0], normFace, 0, nil
}

// 模型热身
func warmup(path string){
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Printf("warmup fail: %s", err.Error())
		return
	}

	for _, file := range files {
		if file.IsDir() { continue }
	
		img, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", path, file.Name()))
		if err != nil { continue }

		r, r2, _, _, err := featuresInfer(img)
		if err==nil {
			log.Printf("warmup: %s %v %v", file.Name(), len(r), len(r2))
		} else {
			log.Printf("warmup fail: %s %s", file.Name(), err.Error())
		}
	}
}

// 1:N
func search_1_N(requestId, groupId string, img []byte) (*map[string]interface{}, error) {
	// 模型推理
	feat, box, normFace, code, err := featuresInfer(img)
	if err != nil {
		return &map[string]interface{}{"code":code}, err
	}

	if feat==nil {  // 未检测到人脸
		return &map[string]interface{}{"user_list":[]int{}}, nil
	}

	// 正则化
	feat, err = norm(feat)
	if err != nil {
		return &map[string]interface{}{"code":9005}, err
	}

	r := gosearch.Search(groupId, feat)

	// FAS 检查
	isReal, realScore, err := fas2.FasCheck(normFace)
	if err != nil {
		return &map[string]interface{}{"code":9007}, err
	}

	// 保存请求图片和结果
	saveBackLog(requestId, img, []byte(fmt.Sprintf("%v %v %v", r, isReal, realScore)))

	if r==nil { // 未识别到 label
		return &map[string]interface{}{"user_list":[]int{}}, nil
	}

	// 检测数据库连接
	if !facelib.Ping() {
		return &map[string]interface{}{"code":9008}, fmt.Errorf("DB connection problem.")
	}

	// 获取用户信息
	database := facelib.Client.Database("face_db")
	collUsers := database.Collection("users")

	var result bson.M
	var opt options.FindOneOptions
	opt.SetProjection(bson.M{"mobile":1, "name":1, "_id":0})

	err = collUsers.FindOne(context.Background(), bson.D{
		{"group_id", groupId}, 
		{"user_id", (*r)["label"].(string) },
	}, &opt).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &map[string]interface{}{"code":9006}, fmt.Errorf("user_id not found")
		} else {
			return &map[string]interface{}{"code":9009}, err
		}
	}

	//log.Printf("%v", result)

	// 取电话号码后4位
	mtail := result["mobile"].(string)
	if len(mtail)>4 {
		mtail = mtail[len(mtail)-4:]
	}

	// 返回结果
	return &map[string]interface{}{ "user_list" : []map[string]interface{}{
			map[string]interface{}{
				"user_id":     (*r)["label"].(string),
				"mobile_tail": mtail,
				"name":        result["name"].(string),
				"location":    box[:4],
				"score":       (*r)["score"].(float32) /2 + 0.5, // 结果 在 [0,1] 之间
				"fake":        []interface{}{!isReal, realScore},
			},
		},
	}, nil
}

// 1:1
func search_1_1(requestId, groupId, userId string, img []byte) (*map[string]interface{}, error) {
	// 模型推理
	feat, box, _, code, err := featuresInfer(img)
	if err != nil {
		return &map[string]interface{}{"code":code}, err
	}

	if feat==nil {  // 未检测到人脸
		return &map[string]interface{}{"is_match":false, "score":0}, nil
	}

	// 正则化
	feat, err = norm(feat)
	if err != nil {
		return &map[string]interface{}{"code":9005}, err
	}

	// 检测数据库连接
	if !facelib.Ping() {
		return &map[string]interface{}{"code":9008}, fmt.Errorf("DB connection problem.")
	}

	// 获取用户信息
	database := facelib.Client.Database("face_db")
	collUsers := database.Collection("users")
	collFaces := database.Collection("faces")

	var user bson.M
	var opt options.FindOneOptions
	opt.SetProjection(bson.M{"face_list":1, "_id":0})

	err = collUsers.FindOne(context.Background(), bson.D{
		{"group_id", groupId}, 
		{"user_id", userId },
	}, &opt).Decode(&user)
	if err != nil { 
		if err == mongo.ErrNoDocuments {
			return &map[string]interface{}{"code":9006}, fmt.Errorf("user_id not found")
		} else {
			return &map[string]interface{}{"code":9009}, err
		}
	}

	// 匹配 face_list
	_, score := is_match(collFaces, feat, user["face_list"].(primitive.A))

	if score<gosearch.ThreshHold {
		// 匹配到
		return &map[string]interface{}{
			"is_match" : true,
			"score"    : score / 2 + 0.5,
			"location" : box[:4],
		}, nil
	} else {
		// 未匹配到
		return &map[string]interface{}{"is_match":false, "score":0}, nil
	}
}


// 双因素：人脸 + 号码厚4位
func search_1_mobile(requestId, groupId, mobileTail string, img []byte) (*map[string]interface{}, error) {
	// 模型推理
	feat, box, _, code, err := featuresInfer(img)
	if err != nil {
		return &map[string]interface{}{"code":code}, err
	}

	if feat==nil {  // 未检测到人脸
		return &map[string]interface{}{"is_match":false, "score":0}, nil
	}

	// 正则化
	feat, err = norm(feat)
	if err != nil {
		return &map[string]interface{}{"code":9005}, err
	}

	// 检测数据库连接
	if !facelib.Ping() {
		return &map[string]interface{}{"code":9008}, fmt.Errorf("DB connection problem.")
	}

	// 获取用户信息
	database := facelib.Client.Database("face_db")
	collUsers := database.Collection("users")
	collFaces := database.Collection("faces")

	//var user bson.M
	var opt options.FindOptions
	opt.SetProjection(bson.M{"face_list":1, "user_id":1, "name":1, "_id":0})

	// 读取用户列表
	cur, err := collUsers.Find(context.Background(), bson.D{
		{"group_id", groupId},
		{"mobile", bson.D{ {"$regex", mobileTail+"$"} } },
	}, &opt)
	if err != nil { 
		return &map[string]interface{}{"code":9009}, err
	}
	defer cur.Close(context.Background())

	// 人脸的 _id 列表
	var faceList primitive.A
	var labelName []string
	var userName []string

	// 输出数据
	for cur.Next(context.Background()) {
		var user bson.M
		if err = cur.Decode(&user); err != nil {
			log.Println("Fetch group data fail: ", err)
			continue
		}

		for _, item := range user["face_list"].(primitive.A){
			labelName = append(labelName, user["user_id"].(string))
			userName = append(labelName, user["name"].(string))
			faceList = append(faceList, item)
		}
	}

	// 匹配 face_list
	pos, score := is_match(collFaces, feat, faceList)

	if score<gosearch.ThreshHold {
		// 匹配到
		return &map[string]interface{}{ "user_list" : []map[string]interface{}{
				map[string]interface{}{
					"user_id"  : labelName[pos],
					"name"     : userName[pos],
					"score"    : score / 2 + 0.5,
					"location" : box[:4],
				},
			},
		}, nil
	} else {
		// 未匹配到
		return &map[string]interface{}{"user_list":[]int{}}, nil
	}
}


// 在 face_list 里比较，匹配最近的一个（score最小的）
func is_match(collFaces *mongo.Collection, feat []float32, face_list primitive.A) (pos int, best float32){
	var opt options.FindOneOptions
	opt.SetProjection(bson.M{"encodings":1, "_id":0})

	best = 1 // 不匹配

	for n, item := range face_list {
		var result bson.M
		objID, _ := primitive.ObjectIDFromHex(item.(string))
		err := collFaces.FindOne(context.Background(), bson.M{"_id": objID}, &opt).Decode(&result)
		if err != nil { 
			log.Println("Find face fail: ", err)
			continue
		}
		//log.Println(result["encodings"].(bson.M)["arc"].(bson.M)["None"].(primitive.A))

		// 获取特征
		encodings, ok := result["encodings"].(bson.M)
		if !ok {
			log.Println("encodings err: ", item)
			continue
		}

		vec2 := make([]float32, vecLen)

		// 装入 arcface 特征
		arc, ok := encodings["arc"].(bson.M)
		if !ok {
			log.Println("arc encodings err: ", item)
			continue
		}

		// 只取 None 人脸特征
		arcNon, ok := arc["None"].(primitive.A)
		if !ok {
			log.Println("arcNone encodings err: ", item)
			continue
		}

		for i := range arcNon {
			vec2[i] = float32(arcNon[i].(float64))
		}

		// 计算 余弦相似度
		score, err := cosine(feat, vec2)
		if err != nil {
			log.Println("calc cosine fail: ", err)
			continue
		}

		if float32(score)<gosearch.ThreshHold {
			// 匹配到
			pos = n
			best = float32(score)
		}
	}

	return
}