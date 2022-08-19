package main

import (
	//"log"
	"fmt"
	"image"

	//"github.com/disintegration/imaging"

	"onnx_test/gocvx"
)

const (
	face_align_image_size = 112
)

var (
	arcface_src = []gocvx.Point2f{
	   {38.2946, 51.6963},
       {73.5318, 51.5014},
       {56.0252, 71.7366},
       {41.5493, 92.3655},
       {70.7299, 92.2041},
   }
)

func norm_crop(srcImage image.Image, lmk []float32) (image.Image, error) {
	// 仿射变化
	m := estimate_norm(lmk)
	defer m.Close()

	// 仿射变化矩阵
	//printM(m)

	// 转换为 Mat
	src, err := gocvx.ImageToMatRGB(srcImage)
	if err!=nil {
		return nil, err
	}

	//log.Println(src.Cols(), src.Rows())

	dst := src.Clone()
	defer dst.Close()

	// 扣图
	gocvx.WarpAffine(src, &dst, m, image.Point{face_align_image_size, face_align_image_size})

	//log.Println(dst.Cols(), dst.Rows())

	aimg, err := dst.ToImage()
	if err!=nil {
		return nil, err
	}

	return aimg, nil
}


// 计算放射 矩阵, 等效 SimilarityTransform()
func estimate_norm(lmk []float32) gocvx.Mat {
	dst := make([]gocvx.Point2f, 5)
	for i:=0;i<5;i++ {
		dst[i] = gocvx.Point2f{lmk[i*2], lmk[i*2+1]}
	}

	pvsrc := gocvx.NewPoint2fVectorFromPoints(arcface_src)
	defer pvsrc.Close()

	pvdst := gocvx.NewPoint2fVectorFromPoints(dst)
	defer pvdst.Close()

	//log.Println(pvdst.ToPoints())
	//log.Println(pvsrc.ToPoints())

	inliers := gocvx.NewMat()
	defer inliers.Close()
	method := 4 // cv2.LMEDS
	ransacProjThreshold := 3.0
	maxiters := uint(2000)
	confidence := 0.99
	refineIters := uint(10)

	m := gocvx.EstimateAffinePartial2DWithParams(pvdst, pvsrc, inliers, method, 
												 ransacProjThreshold, maxiters, confidence, refineIters)
	//defer m.Close()

	return m
}


func printM(m gocvx.Mat) {
	for i:=0;i<m.Rows();i++ {
		for j:=0;j<m.Cols();j++ {
			fmt.Printf("%v ", m.GetDoubleAt(i, j))
		}
		fmt.Printf("\n")
	}
}