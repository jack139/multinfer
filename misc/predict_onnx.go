package main

/*
CGO_LDFLAGS="-L/usr/local/lib -lopencv_core -lopencv_calib3d -lopencv_imgproc" go build
LD_LIBRARY_PATH=/usr/local/lib ./onnx_test
*/

import (
	"log"
	"image"

	"github.com/ivansuteja96/go-onnxruntime"
	"github.com/disintegration/imaging"

	"onnx_test/gocvx"
)

const (
	test_image_path = "data/5.jpg"
	detModel_path = "../../../cv/face_model/arcface/models/buffalo_l/det_10g.onnx"

	det_model_input_size = 224
)

func main() {
	ortEnvDet := onnxruntime.NewORTEnv(onnxruntime.ORT_LOGGING_LEVEL_WARNING, "development")
	ortDetSO := onnxruntime.NewORTSessionOptions()

	detModel, err := onnxruntime.NewORTSession(ortEnvDet, detModel_path, ortDetSO)
	if err != nil {
		log.Println(err)
		return
	}

	// load image
	srcImage, err := imaging.Open(test_image_path)
	if err != nil {
		log.Fatal("Error: %s\n", err.Error())
	}

	shape1 := []int64{1, 3, det_model_input_size, det_model_input_size}
	input1, det_scale := preprocessImage(srcImage, det_model_input_size)

	//log.Println(input1[:100])

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

	log.Println(dets)
	log.Println(kpss)

	m := estimate_norm(kpss[0])
	defer m.Close()

	printM(m)

	//estimate_affine()

	src, _ := gocvx.ImageToMatRGB(srcImage)
	log.Println(src.Cols(), src.Rows())

	dst := src.Clone()
	defer dst.Close()

	gocvx.WarpAffine(src, &dst, m, image.Point{112, 112})

	log.Println(dst.Cols(), dst.Rows())

	aimg, _ := dst.ToImage()

	_ = imaging.Save(aimg, "data/aimg.jpg")
}



