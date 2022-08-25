package facelib

import (
	"context"
	"log"
	"errors"
	"strings"
	"io/ioutil"
	"strconv"
	//"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	vggLen = 2048
	evoLen = 512
	arcLen = 512
	//vecLen = vggLen + evoLen
	vecLen = arcLen

	USEARCFACE = true
)

var (
	LimitFace = int(3)  // 注册的人脸数

	X map[string][][]float32 // 特征向量
	y map[string][]uint32     // 标签
	labelName map[string][]string  // 标签对应用户id

	//N int = 0 // 脸特征总数
)

func ReadData(groupStr string) error {
	log.Printf("GONUM= %d\tLimitFace= %d\n", GONUM, LimitFace)

	if !Ping() {
		return errors.New("DB connection problem.")
	}

	/*
		取特征数据
		1. 去指定group的用户
		2. 每个用户，取人脸特征
	*/
	database := Client.Database("face_db")
	collUsers := database.Collection("users")
	collFaces := database.Collection("faces")

	// 用户组列表
	groupList := strings.Split(groupStr, ",")

	X = make(map[string][][]float32)
	y = make(map[string][]uint32)
	labelName = make(map[string][]string)

	// 分组读入特征数据
	for _, group := range groupList {

		// 读取用户列表
		cur, err := collUsers.Find(context.Background(), bson.D{{"group_id", group}})
		if err != nil { 
			log.Println("Find group fail: ", err)
			continue
		}
		defer cur.Close(context.Background())

		// 人脸的 _id 列表
		var faceList []primitive.A

		// 输出数据
		for cur.Next(context.Background()) {
			var user bson.M
			if err = cur.Decode(&user); err != nil {
				log.Println("Fetch group data fail: ", err)
				continue
			}

			labelName[group] = append(labelName[group], user["user_id"].(string))
			faceList = append(faceList, user["face_list"].(primitive.A))
		}

		// 读取每个用户的 人脸特征
		var opt options.FindOneOptions
		opt.SetProjection(bson.M{"encodings": 1})

		for label, _ := range labelName[group] {

			// 读取特征数据
			for n, item := range faceList[label] {
				//log.Println(n, item)
				if n==LimitFace { // 只记录指定数量的特征
					break
				}
				var result bson.M
				objID, _ := primitive.ObjectIDFromHex(item.(string))
				err := collFaces.FindOne(context.Background(), bson.M{"_id": objID}, &opt).Decode(&result)
				if err != nil { 
					log.Println("Find face fail: ", err)
					continue
				}
				//log.Println(result["encodings"].(bson.M)["evo"].(bson.M)["360"].(primitive.A))

				// 获取特征
				encodings, ok := result["encodings"].(bson.M)
				if !ok {
					log.Println("encodings err: ", item)
					continue
				}

				vec2 := make([]float32, vecLen)

				if USEARCFACE {
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
				} else { 
					// 装入 vgg+evo 特征
					vgg, ok := encodings["vgg"].(bson.M)
					if !ok {
						log.Println("vgg encodings err: ", item)
						continue
					}
					evo, ok := encodings["evo"].(bson.M)
					if !ok {
						log.Println("evo encodings err: ", item)
						continue
					}

					// 只取 None 人脸特征
					vggNon, ok := vgg["None"].(primitive.A)
					if !ok {
						log.Println("vggNone encodings err: ", item)
						continue
					}
					evoNon, ok := evo["None"].(primitive.A)
					if !ok {
						log.Println("evoNone encodings err: ", item)
						continue
					}

					for i := range vggNon {
						vec2[i] = float32(vggNon[i].(float64))
					}
					for i := range evoNon {
						vec2[i+vggLen] = float32(evoNon[i].(float64))
					}
				}

				X[group] = append(X[group], vec2)
				y[group] = append(y[group], uint32(label+1))

			}

		}

		log.Printf("%s: label_n= %d\tX_n= %d\ty_n= %d\n", group, len(labelName[group]), len(X[group]), len(y[group]))
		if len(X[group])>0 {
			log.Printf("\tdim= %d\n", len(X[group][0]))
		}
	}

	return nil
}

// 新增特征，只在内存新增，不处理数据库
func AddNewData(group string, label string, vec []float32){
	labelName[group] = append(labelName[group], label)
	X[group] = append(X[group], vec)
	y[group] = append(y[group], uint32(len(labelName[group])))

	//log.Println(labelName[group])
	log.Println("Add data: ", group, label)
}

// 删除特征，不处理数据库，只删除labelname，不真的删除特征
func RemoveData(group string, label string){
	for n, k := range labelName[group] {
		if k==label {
			labelName[group][n] = "__BLANK__"
		}
	}

	//log.Println(labelName[group])
	log.Println("Remove data: ", group, label)
}


func searchLabel(label string, group string) int {
	for n, v := range labelName[group] {
		if v==label {
			return n
		}
	}
	return -1
}

// 从文件载入注册数据 (测试用)
func ReadDataFromFile(dataFile string, group string){
	log.Printf("GONUM= %d\tLimitFace= %d\n", GONUM, LimitFace)

	X = make(map[string][][]float32)
	y = make(map[string][]uint32)
	labelName = make(map[string][]string)

	b, err := ioutil.ReadFile(dataFile) 
	if err != nil {
		log.Println(err)
		return
	}
	s := string(b)
	lines := strings.Split(s, "\n")

	//log.Println(len(lines), len(lines[0]))

	faceNum := 0 // 记录同一label的脸数量

	for i:=0;i<len(lines);i++ {
		if len(lines[i])==0 { continue } // 过滤掉空行
		xx := strings.Split(lines[i], ",")
		
		// 取得一个人脸数据
		var pos int
		vec := make([]float32, 0)
		for xn,fs := range xx {
			if xn==0 { // 第一个字段是 label
				pos = searchLabel(fs, group)
				if pos==-1 { // 这里假设同一人人脸数据是连续的
					labelName[group] = append(labelName[group], fs)
					pos = len(labelName[group])-1	
					faceNum = 0
				}
				continue
			}
			f, _ := strconv.ParseFloat(fs, 32)
			vec = append(vec, float32(f))
			//log.Printf("%.8f ", f)
		}

		if faceNum==LimitFace { // 只记录指定数量的特征
			continue
		}

		y[group] = append(y[group], uint32(pos+1))
		X[group] = append(X[group], vec)

		faceNum += 1

		//log.Println()
	}

	log.Printf("%s: label_n= %d\tX_n= %d\ty_n= %d\n", group, len(labelName[group]), len(X[group]), len(y[group]))
	if len(X[group])>0 {
		log.Printf("\tdim= %d\n", len(X[group][0]))
	}
}
