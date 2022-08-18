package main

import (
	"log"
	"fmt"
	//"image"

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

	log.Println(pvdst.ToPoints())
	log.Println(pvsrc.ToPoints())

	inliers := gocvx.NewMat()
	defer inliers.Close()
	method := 4 // cv2.LMEDS
	ransacProjThreshold := 3.0
	maxiters := uint(2000)
	confidence := 0.99
	refineIters := uint(10)

	m := gocvx.EstimateAffinePartial2DWithParams(pvdst, pvsrc, inliers, method, ransacProjThreshold, maxiters, confidence, refineIters)
	//defer m.Close()

	return m
}

// 测试
func estimate_affine() {
	dst := []gocvx.Point2f{
		{218.78867, 205.74413},
		{312.13818, 202.18082},
		{279.89087, 232.69415},
		{236.05072, 302.79538},
		{313.98624, 299.34445},
	}

	src := []gocvx.Point2f{
		{38.2946, 51.6963},
		{73.5318, 51.5014},
		{56.0252, 71.7366},
		{41.5493, 92.3655},
		{70.7299, 92.2041},
	}


	pvsrc := gocvx.NewPoint2fVectorFromPoints(src)
	defer pvsrc.Close()

	pvdst := gocvx.NewPoint2fVectorFromPoints(dst)
	defer pvdst.Close()

	log.Println(pvdst.ToPoints())
	log.Println(pvsrc.ToPoints())

	inliers := gocvx.NewMat()
	defer inliers.Close()
	method := 4 // cv2.LMEDS
	ransacProjThreshold := 3.0
	maxiters := uint(2000)
	confidence := 0.99
	refineIters := uint(10)

	m := gocvx.EstimateAffinePartial2DWithParams(pvdst, pvsrc, inliers, method, ransacProjThreshold, maxiters, confidence, refineIters)
	//m := EstimateAffinePartial2D(pvdst, pvsrc)
	defer m.Close()

	log.Println(m.Type(), m.Step())

	printM(m)
	//printM(inliers)

	//v, _ := m.DataPtrFloat64()
	//log.Println(v)	
}

func printM(m gocvx.Mat) {
	for i:=0;i<m.Rows();i++ {
		for j:=0;j<m.Cols();j++ {
			fmt.Printf("%v ", m.GetDoubleAt(i, j))
		}
		fmt.Printf("\n")
	}
}