#include <string>
#include <iostream>

#include "predict.h"

Mat EstimateAffine2DWithParams(Point2fVector from, Point2fVector to, Mat inliers, int method, double ransacReprojThreshold, size_t maxIters, double confidence, size_t refineIters) {
	std::cerr << "----test\n";
    return new cv::Mat(cv::estimateAffine2D(*from, *to, *inliers, method, ransacReprojThreshold, maxIters, confidence, refineIters));
}
