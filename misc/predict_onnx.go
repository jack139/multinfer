package main

import (
	"fmt"
	"os"
	"image"
	"image/color"
	"log"

	"github.com/ivansuteja96/go-onnxruntime"
	"github.com/disintegration/imaging"
)

// LD_LIBRARY_PATH=/usr/local/lib go run predict_example2.go
func main() {
	ortEnvDet := onnxruntime.NewORTEnv(onnxruntime.ORT_LOGGING_LEVEL_VERBOSE, "development")
	ortDetSO := onnxruntime.NewORTSessionOptions()

	detModel, err := onnxruntime.NewORTSession(ortEnvDet, "../../multinfer/data/det_10g.onnx", ortDetSO)
	if err != nil {
		log.Println(err)
		return
	}

	shape1 := []int64{1, 3, 224, 224}
	input1 := preprocessImage("../../source/5a.jpg", 224)

	fmt.Println(input1[:100])

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

	for i:=0;i<len(res);i++ {
		fmt.Printf("Success do predict, shape : %+v, result : %+v\n", 
			res[i].Shape, 
			res[i].Value.([]float32)[:res[i].Shape[1]], // only show one value
		)
	}
}


func Transpose(rgbs []float32) []float32 {
	out := make([]float32, len(rgbs))
	channelLength := len(rgbs) / 3
	for i := 0; i < channelLength; i++ {
		out[i] = rgbs[i*3]
		out[i+channelLength] = rgbs[i*3+1]
		out[i+channelLength*2] = rgbs[i*3+2]

		// RGB --> BGR
		//out[i] = rgbs[i*3+2]
		//out[i+channelLength] = rgbs[i*3+1]
		//out[i+channelLength*2] = rgbs[i*3]
	}
	return out
}

func preprocessImage(imageFile string, inputSize int) []float32 {
	src, err := imaging.Open(imageFile)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	rgbs := make([]float32, inputSize*inputSize*3)

	result := imaging.Resize(src, 163, 224, imaging.Lanczos)
	fmt.Println("resize: ", result.Rect)
	//result = imaging.CropAnchor(result, 224, 224, imaging.Center)
	//fmt.Println("crop: ", result.Rect)
	result = padBox(result)
	_ = imaging.Save(result, "/tmp/test2.jpg")

	j := 0
	for i := range result.Pix {
		if (i+1)%4 != 0 {
			rgbs[j] = float32(result.Pix[i])
			j++
		}
	}

	fmt.Println(rgbs[:100])

	rgbs = Transpose(rgbs)

	fmt.Println(rgbs[:100])

	channelLength := len(rgbs) / 3
	for i := 0; i < channelLength; i++ {
		rgbs[i] = normalize(rgbs[i], 127.5, 128.0)
		rgbs[i+channelLength] = normalize(rgbs[i+channelLength], 127.5, 128.0)
		rgbs[i+channelLength*2] = normalize(rgbs[i+channelLength*2], 127.5, 128.0)
	}
	return rgbs
}

func normalize(in float32, m float32, s float32) float32 {
	return (in - m) / s
}


// 调整为方形，黑色填充
func padBox(src image.Image) *image.NRGBA {
	var maxW int

	if src.Bounds().Dx() > src.Bounds().Dy() {
		maxW = src.Bounds().Dx()
	} else {
		maxW = src.Bounds().Dy()
	}

	dst := imaging.New(maxW, maxW, color.Black)
	dst = imaging.Paste(dst, src, image.Point{0,0})

	//_ = imaging.Save(dst, "data/test3.jpg")

	return dst
}
