package fas2

import (
	"log"
	"math"
	"image"
	"image/color"

	"github.com/disintegration/imaging"
)


func transposeRGB(rgbs []float32) []float32 {
	out := make([]float32, len(rgbs))
	channelLength := len(rgbs) / 3
	for i := 0; i < channelLength; i++ {
		//out[i] = rgbs[i*3]
		//out[i+channelLength] = rgbs[i*3+1]
		//out[i+channelLength*2] = rgbs[i*3+2]

		// RGB --> BGR
		out[i] = rgbs[i*3+2]
		out[i+channelLength] = rgbs[i*3+1]
		out[i+channelLength*2] = rgbs[i*3]
	}
	return out
}

func preprocessImage(src image.Image, inputSize int) []float32 {
	var newHeight, newWidth int
	im_ratio := float32(src.Bounds().Dx()) / float32(src.Bounds().Dy())
	if im_ratio > 1 { // width > height
		newWidth = inputSize
		newHeight = int(float32(newWidth) / im_ratio)
	} else {
		newHeight = inputSize
		newWidth = int(float32(newHeight) * im_ratio)		
	}

	//result := imaging.Clone(src)
	//log.Println(src.Bounds(), newWidth, newHeight)
	//result := imaging.Resize(src, inputSize, inputSize, imaging.Lanczos) // 直接调整尺寸，会变形
	result := imaging.Resize(src, newWidth, newHeight, imaging.Lanczos)	// 保持比例
	//log.Println("resize: ", result.Rect)
	result = padBox(result)


	rgbs := make([]float32, inputSize*inputSize*3)

	j := 0
	for i := range result.Pix {
		if (i+1)%4 != 0 {
			rgbs[j] = float32(result.Pix[i])
			j++
		}
	}

	//log.Println(rgbs[:100])

	rgbs = transposeRGB(rgbs)

	//log.Println(rgbs[:100])

	//channelLength := len(rgbs) / 3
	//for i := 0; i < channelLength; i++ {
	//	rgbs[i] = normalize(rgbs[i], 127.5, 128.0)
	//	rgbs[i+channelLength] = normalize(rgbs[i+channelLength], 127.5, 128.0)
	//	rgbs[i+channelLength*2] = normalize(rgbs[i+channelLength*2], 127.5, 128.0)
	//}

	//log.Println("det_scale===", det_scale, float32(newHeight), float32(src.Bounds().Dy()))

	return rgbs
}

func normalize(in float32, m float32, s float32) float32 {
	return (in - m) / s
}


// 调整为方形，黑色填充, 图片居中
func padBox(src image.Image) *image.NRGBA {
	var maxW int

	if src.Bounds().Dx() > src.Bounds().Dy() {
		maxW = src.Bounds().Dx()
	} else {
		maxW = src.Bounds().Dy()
	}

	dst := imaging.New(maxW, maxW, color.Black)
	//dst = imaging.Paste(dst, src, image.Point{0,0})
	dst = imaging.PasteCenter(dst, src)

	//_ = imaging.Save(dst, "data/test2.jpg")

	return dst
}


func softmax(x []float32) []float32 {
	var sum float64
	x2 := make([]float64, len(x))
	y := make([]float32, len(x))

	for i:=0;i<len(x);i++ {
		x2[i] = math.Exp(float64(x[i]))
		sum += x2[i]
	}

	log.Println("sum=", sum)

	for i:=0;i<len(x);i++ {
		y[i] = float32(x2[i] / sum)
	}

	return y
}