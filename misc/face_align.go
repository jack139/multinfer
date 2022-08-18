package main

import (
	"log"
	"fmt"
	//"image"

	//"github.com/disintegration/imaging"

	"onnx_test/gocv"
)

func estimate_affine() {
	dst := []gocv.Point2f{
		{218.78867, 205.74413},
		{312.13818, 202.18082},
		{279.89087, 232.69415},
		{236.05072, 302.79538},
		{313.98624, 299.34445},
	}

	src := []gocv.Point2f{
		{38.2946, 51.6963},
		{73.5318, 51.5014},
		{56.0252, 71.7366},
		{41.5493, 92.3655},
		{70.7299, 92.2041},
	}


	pvsrc := gocv.NewPoint2fVectorFromPoints(src)
	defer pvsrc.Close()

	pvdst := gocv.NewPoint2fVectorFromPoints(dst)
	defer pvdst.Close()

	log.Println(pvdst.ToPoints())
	log.Println(pvsrc.ToPoints())

	inliers := gocv.NewMat()
	defer inliers.Close()
	method := 4 // cv2.LMEDS
	ransacProjThreshold := 3.0
	maxiters := uint(2000)
	confidence := 0.99
	refineIters := uint(10)

	m := gocv.EstimateAffinePartial2DWithParams(pvdst, pvsrc, inliers, method, ransacProjThreshold, maxiters, confidence, refineIters)
	//m := EstimateAffinePartial2D(pvdst, pvsrc)
	defer m.Close()

	printM(m)
	printM(inliers)

	log.Println(m.Type(), m.Step())

	v, _ := m.DataPtrFloat64()
	log.Println(v)	
}

func printM(m gocv.Mat) {
	for i:=0;i<m.Rows();i++ {
		for j:=0;j<m.Cols();j++ {
			fmt.Printf("%v ", m.GetDoubleAt(i, j))
		}
		fmt.Printf("\n")
	}
}