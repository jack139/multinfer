package gocvx

/*
#include <stdlib.h>
#include "gocvx.h"
*/
import "C"
import (
	"image"
	//"image/color"
	//"errors"
)

// EstimateAffinePartial2D computes an optimal limited affine transformation
// with 4 degrees of freedom between two 2D point sets.
//
// For further details, please see:
// https://docs.opencv.org/master/d9/d0c/group__calib3d.html#gad767faff73e9cbd8b9d92b955b50062d
//
// add more parameters to original gocv EstimateAffinePartial2D()
func EstimateAffinePartial2DWithParams(from Point2fVector, to Point2fVector, inliers Mat, method int, ransacReprojThreshold float64, maxIters uint, confidence float64, refineIters uint) Mat {
	return newMat(C.EstimateAffinePartial2DWithParams(from.p, to.p, inliers.p, C.int(method), C.double(ransacReprojThreshold), C.size_t(maxIters), C.double(confidence), C.size_t(refineIters)))
}


// WarpAffine applies an affine transformation to an image. For more parameters please check WarpAffineWithParams
//
// For further details, please see:
// https://docs.opencv.org/master/da/d54/group__imgproc__transform.html#ga0203d9ee5fcd28d40dbc4a1ea4451983
func WarpAffine(src Mat, dst *Mat, m Mat, sz image.Point) {
	pSize := C.struct_Size{
		width:  C.int(sz.X),
		height: C.int(sz.Y),
	}

	C.WarpAffine(src.p, dst.p, m.p, pSize)
}

/*
// ImageToMatRGB converts image.Image to gocv.Mat,
// which represents RGB image having 8bit for each component.
// Type of Mat is gocv.MatTypeCV8UC3.
func ImageToMatRGB(img image.Image) (Mat, error) {
	bounds := img.Bounds()
	x := bounds.Dx()
	y := bounds.Dy()

	var data []uint8
	switch img.ColorModel() {
	case color.RGBAModel:
		m, res := img.(*image.RGBA)
		if true != res {
			return NewMat(), errors.New("Image color format error")
		}
		data = m.Pix
		// speed up the conversion process of RGBA format
		src, err := NewMatFromBytes(y, x, MatTypeCV8UC4, data)
		if err != nil {
			return NewMat(), err
		}
		defer src.Close()

		dst := NewMat()
		CvtColor(src, &dst, ColorRGBAToBGR)
		return dst, nil

	default:
		data := make([]byte, 0, x*y*3)
		for j := bounds.Min.Y; j < bounds.Max.Y; j++ {
			for i := bounds.Min.X; i < bounds.Max.X; i++ {
				r, g, b, _ := img.At(i, j).RGBA()
				data = append(data, byte(b>>8), byte(g>>8), byte(r>>8))
			}
		}
		return NewMatFromBytes(y, x, MatTypeCV8UC3, data)
	}
}
*/