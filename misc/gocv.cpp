#include <string>
#include <iostream>
#include <vector>

#include "gocv.h"

/*
// for test: print Point2fVector contents
void printVVVVV(Point2fVector const &input)
{
	for (int i = 0; i < input->size(); i++) {
		std::cerr << input->at(i) << ' ';
	}
	std::cerr << "\n";
}
*/

Mat EstimateAffinePartial2DWithParams(Point2fVector from, Point2fVector to, Mat inliers, int method, double ransacReprojThreshold, size_t maxIters, double confidence, size_t refineIters) {
	//printVVVVV(from);
	//printVVVVV(to);
	return new cv::Mat(cv::estimateAffinePartial2D(*from, *to, *inliers, method, ransacReprojThreshold, maxIters, confidence, refineIters));
}

/*
// declaration in opencv 4.5.5
CV_EXPORTS_W cv::Mat estimateAffinePartial2D(InputArray from, InputArray to, OutputArray inliers = noArray(),
								  int method = RANSAC, double ransacReprojThreshold = 3,
								  size_t maxIters = 2000, double confidence = 0.99,
								  size_t refineIters = 10);
*/