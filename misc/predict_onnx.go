package main

import (
	"log"

	"github.com/ivansuteja96/go-onnxruntime"
)

const (
	det_model_input_size = 224
	nms_thresh = float32(0.4)
	det_thresh = float32(0.5)
)

// CGO_LDFLAGS="-L/usr/local/lib -lopencv_core -lopencv_calib3d -lopencv_imgproc" go build
// LD_LIBRARY_PATH=/usr/local/lib ./onnx_test
func main() {
	ortEnvDet := onnxruntime.NewORTEnv(onnxruntime.ORT_LOGGING_LEVEL_WARNING, "development")
	ortDetSO := onnxruntime.NewORTSessionOptions()

	detModel, err := onnxruntime.NewORTSession(ortEnvDet, "../../../cv/face_model/arcface/models/buffalo_l/det_10g.onnx", ortDetSO)
	if err != nil {
		log.Println(err)
		return
	}

	shape1 := []int64{1, 3, det_model_input_size, det_model_input_size}
	input1, det_scale := preprocessImage("data/5.jpg", det_model_input_size)

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

	estimate_affine()
}



