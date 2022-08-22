package yhface

import (
	"log"
	"fmt"
	"bytes"
	"io/ioutil"

	"github.com/disintegration/imaging"
	"github.com/jack139/go-infer/helper"
	"github.com/jack139/arcface-go/arcface"
)

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

		log.Println("Onnx model loaded from: ", helper.Settings.Customer["ArcfaceModelPath"])

		initOK = true

		// 模型热身
		//warmup(helper.Settings.Customer["FACE_WARM_UP_IMAGES"])
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


func featuresInfer(imageByte []byte) ([]float32, int, error){

	// 转换为 image.Image
	reader := bytes.NewReader(imageByte)

	img, err := imaging.Decode(reader)
	if err!=nil {
		return nil, 9201,err
	}

	// 检测人脸
	dets, kpss, err := arcface.FaceDetect(img)
	if err != nil {
		return nil, 9202, err
	}

	if len(dets)==0 {
		log.Println("No face detected.")
		return nil, 0, nil
	}

	// 只返回第一个人脸的特征
	features, err := arcface.FaceFeatures(img, kpss[0])
	if err != nil {
		return nil, 9203, err
	}

	return features, 0, nil
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
	
		image, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", path, file.Name()))
		if err != nil { continue }

		r, _, err := locateInfer(image)
		if err==nil {
			log.Printf("warmup: %s %s", file.Name(), r)
		} else {
			log.Printf("warmup fail: %s %s", file.Name(), err.Error())
		}
	}
}
