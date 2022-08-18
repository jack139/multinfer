package main

/*
#include <stdlib.h>
#include "gocv.h"
*/
import "C"


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