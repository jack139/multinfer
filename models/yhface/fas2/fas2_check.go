package fas2

import (
	"log"
	"image"
	"errors"
	"path/filepath"

	"github.com/ivansuteja96/go-onnxruntime"
)

const (
	fas_model_input_size = 80
)

var (
	fas2Model1 *onnxruntime.ORTSession
	fas2Model2 *onnxruntime.ORTSession
)



func FasCheck(src image.Image) (live bool, score float32, err error) {
	shape1 := []int64{1, 3, fas_model_input_size, fas_model_input_size}
	input1 := preprocessImage(src, fas_model_input_size)

	//log.Println(input1[:100])

	// 模型1 检测
	res, err := fas2Model1.Predict([]onnxruntime.TensorValue{
		{
			Value: input1,
			Shape: shape1,
		},
	})
	if err != nil {
		return
	}

	if len(res) == 0 {
		err = errors.New("Fail to get result")
		return
	}

	predictionA1 := res[0].Value.([]float32)
	predictionB1 := softmax(predictionA1)

	log.Println("predictionA1:", predictionA1)
	log.Println("predictionB1:", predictionB1)



	// 模型2 检测
	res2, err := fas2Model2.Predict([]onnxruntime.TensorValue{
		{
			Value: input1,
			Shape: shape1,
		},
	})
	if err != nil {
		return
	}

	if len(res2) == 0 {
		err = errors.New("Fail to get result")
		return
	}

	predictionA2 := res2[0].Value.([]float32)
	predictionB2 := softmax(predictionA2)

	log.Println("predictionA2:", predictionA2)
	log.Println("predictionB2:", predictionB2)

	predictionB1[0] += predictionB2[0]
	predictionB1[1] += predictionB2[1]
	predictionB1[2] += predictionB2[2]

	log.Println("Real Score: ", predictionB1 )

	return predictionB1[1]>predictionB1[0] && predictionB1[1]>predictionB1[2], predictionB1[1] / 2, nil
}


func LoadOnnxModel(onnxmodel_path string) (err error) {
	ortEnvDet := onnxruntime.NewORTEnv(onnxruntime.ORT_LOGGING_LEVEL_ERROR, "development")
	ortDetSO := onnxruntime.NewORTSessionOptions()

	fas2Model1, err = onnxruntime.NewORTSession(ortEnvDet, filepath.Join(onnxmodel_path, "2.7_80x80_MiniFASNetV2.onnx"), ortDetSO)
	if err != nil {
		return err
	}

	fas2Model2, err = onnxruntime.NewORTSession(ortEnvDet, filepath.Join(onnxmodel_path, "4_0_0_80x80_MiniFASNetV1SE.onnx"), ortDetSO)
	if err != nil {
		return err
	}

	return nil
}
