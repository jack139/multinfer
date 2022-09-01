package fas2

import (
	"math"
	"image"

	"github.com/disintegration/imaging"
)


func transposeRGB(rgbs []float32) []float32 {
	out := make([]float32, len(rgbs))
	channelLength := len(rgbs) / 3
	for i := 0; i < channelLength; i++ {
		// RGB --> BGR
		out[i] = rgbs[i*3+2]
		out[i+channelLength] = rgbs[i*3+1]
		out[i+channelLength*2] = rgbs[i*3]
	}
	return out
}


func preprocessImage(src image.Image, inputSize int) []float32 {
	// 直接调整尺寸，normFace已经是正方形，所以不会变形
	result := imaging.Resize(src, inputSize, inputSize, imaging.Lanczos) 

	//_ = imaging.Save(result, "data/pad.jpg")

	rgbs := make([]float32, inputSize*inputSize*3)

	j := 0
	for i := range result.Pix {
		if (i+1)%4 != 0 {
			rgbs[j] = float32(result.Pix[i])
			j++
		}
	}

	rgbs = transposeRGB(rgbs)

	return rgbs
}


func softmax(x []float32) []float32 {
	var sum float64
	x2 := make([]float64, len(x))
	y := make([]float32, len(x))

	for i:=0;i<len(x);i++ {
		x2[i] = math.Exp(float64(x[i]))
		sum += x2[i]
	}

	for i:=0;i<len(x);i++ {
		y[i] = float32(x2[i] / sum)
	}

	return y
}