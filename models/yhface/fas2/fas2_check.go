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


// 返回 live == true 真脸 false 假脸， score 真脸的得分
func FasCheck(src image.Image) (live bool, score float32, err error) {
	shape1 := []int64{1, 3, fas_model_input_size, fas_model_input_size}
	input1 := preprocessImage(src, fas_model_input_size)


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

	predictionA := res[0].Value.([]float32)
	predictionB := softmax(predictionA)


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

	predictionA = res2[0].Value.([]float32)
	predictionA = softmax(predictionA)


	predictionB[0] += predictionA[0]
	predictionB[1] += predictionA[1]
	predictionB[2] += predictionA[2]

	log.Println("Real Score: ", predictionB )

	return predictionB[1]>predictionB[0] && predictionB[1]>predictionB[2], predictionB[1] / 2, nil
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
