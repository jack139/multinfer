#ifndef _PREDICT_H_
#define _PREDICT_H_

#ifdef __cplusplus
#include <opencv2/opencv.hpp>
#include <opencv2/calib3d.hpp>


extern "C" {
#endif

#include "core.h"

Mat EstimateAffinePartial2DWithParams(Point2fVector from, Point2fVector to, Mat inliers, int method, double ransacReprojThreshold, size_t maxIters, double confidence, size_t refineIters);
void WarpAffine(Mat src, Mat dst, Mat m, Size dsize);
void CvtColor(Mat src, Mat dst, int code);

#ifdef __cplusplus
}
#endif

#endif //_PREDICT_H_