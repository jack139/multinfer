package main

/*
CGO_LDFLAGS="-L/usr/local/lib -lopencv_core -lopencv_calib3d -lopencv_imgproc" go build -o data/
LD_LIBRARY_PATH=/usr/local/lib data/onnx_test
*/

import (
	"log"
	//"image"

	"github.com/ivansuteja96/go-onnxruntime"
	"github.com/disintegration/imaging"
)

const (
	test_image_path = "data/6.jpg"
	test_aim_path = "data/aimg.jpg"
	detModel_path = "../../../cv/face_model/arcface/models/buffalo_l/det_10g.onnx"
	arcfaceModel_path = "../../../cv/face_model/arcface/models/buffalo_l/w600k_r50.onnx"

	det_model_input_size = 224
)

var (
	detModel *onnxruntime.ORTSession
	arcfaceModel *onnxruntime.ORTSession
)

func main() {
	if err := loadModelWeight(); err!=nil {
		log.Fatal("Load model fail: ", err.Error())
	}

	// load image
	srcImage, err := imaging.Open(test_image_path)
	if err != nil {
		log.Fatal("Error: %s\n", err.Error())
	}

	shape1 := []int64{1, 3, det_model_input_size, det_model_input_size}
	input1, det_scale := preprocessImage(srcImage, det_model_input_size)

	//log.Println(input1[:100])

	// 人脸检测模型
	res, err := detModel.Predict([]onnxruntime.TensorValue{
		{
			Value: input1,
			Shape: shape1,
		},
	})
	if err != nil {
		log.Println(err)
		return
	}

	if len(res) == 0 {
		log.Println("Failed get result")
		return
	}

	dets, kpss := processResult(res, det_scale)

	//log.Println(dets)
	//log.Println(kpss)

	log.Println("face num: ", len(kpss))

	if len(dets)==0 {
		log.Println("No face detected.")
		return		
	}

	// 截取 校正后的人脸, 只取第一个人脸
	aimg, err := norm_crop(srcImage, kpss[0])
	if err!=nil {
		log.Println(err)
		return		
	}

	_ = imaging.Save(aimg, "data/aimg.jpg")


	//// 截取的框， 未校正的人脸
	//sr := image.Rectangle{
	//	image.Point{int(dets[0][0]), int(dets[0][1])}, 
	//	image.Point{int(dets[0][2]), int(dets[0][3])},
	//}
	//// 截取
	//src2 := imaging.Crop(srcImage, sr)
	//_ = imaging.Save(src2, "data/img.jpg")


	// load image -- 测试
	//aimg, err = imaging.Open(test_aim_path)
	//if err != nil {
	//	log.Fatal("Error: %s\n", err.Error())
	//}

	// 准备数据： 人脸特征模型
	shape2 := []int64{1, 3, face_align_image_size, face_align_image_size}
	input2 := preprocessFace(aimg, face_align_image_size)


	// 人脸特征模型
	res2, err := arcfaceModel.Predict([]onnxruntime.TensorValue{
		{
			Value: input2,
			Shape: shape2,
		},
	})
	if err != nil {
		log.Println(err)
		return
	}

	if len(res2) == 0 {
		log.Println("Failed get result")
		return
	}

	features := res2[0].Value.([]float32)

	log.Println("features: ", features)
	log.Println("shape", res2[0].Shape)
}


func loadModelWeight() (err error) {
	ortEnvDet := onnxruntime.NewORTEnv(onnxruntime.ORT_LOGGING_LEVEL_ERROR, "development")
	ortDetSO := onnxruntime.NewORTSessionOptions()

	detModel, err = onnxruntime.NewORTSession(ortEnvDet, detModel_path, ortDetSO)
	if err != nil {
		return err
	}

	arcfaceModel, err = onnxruntime.NewORTSession(ortEnvDet, arcfaceModel_path, ortDetSO)
	if err != nil {
		return err
	}

	return nil
}
