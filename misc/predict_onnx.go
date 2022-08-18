package main


import (
	"fmt"
	"os"
	"image"
	"image/color"
	"log"
	"sort"

	"github.com/ivansuteja96/go-onnxruntime"
	"github.com/disintegration/imaging"
	//"gocv.io/x/gocv"
)

const (
	det_model_input_size = 224
	nms_thresh = float32(0.4)
	det_thresh = float32(0.5)
)

// LD_LIBRARY_PATH=/usr/local/lib go run predict_onnx.go
// CGO_CPPFLAGS="-I/usr/local/include/opencv4" CGO_LDFLAGS="-L/usr/local/lib -lopencv_core -lopencv_calib3d" go build -o predict_onnx
// LD_LIBRARY_PATH=/usr/local/lib ./predict_onnx
func main() {
	ortEnvDet := onnxruntime.NewORTEnv(onnxruntime.ORT_LOGGING_LEVEL_WARNING, "development")
	ortDetSO := onnxruntime.NewORTSessionOptions()

	detModel, err := onnxruntime.NewORTSession(ortEnvDet, "../../../cv/face_model/arcface/models/buffalo_l/det_10g.onnx", ortDetSO)
	if err != nil {
		log.Println(err)
		return
	}

	shape1 := []int64{1, 3, det_model_input_size, det_model_input_size}
	input1, det_scale := preprocessImage("data/5.jpg", det_model_input_size)

	//fmt.Println(input1[:100])

	res, err := detModel.Predict([]onnxruntime.TensorValue{
		{
			Value: input1,
			Shape: shape1,
		},
	})
	if err != nil {
		log.Println(err)
		return
	}

	if len(res) == 0 {
		log.Println("Failed get result")
		return
	}

	dets, kpss := processResult(res, det_scale)

	fmt.Println(dets)
	fmt.Println(kpss)

	//fmt.Printf("gocv version: %s\n", Version())
	//fmt.Printf("opencv lib version: %s\n", OpenCVVersion())

	estimate_affine()
}


func TransposeRGB(rgbs []float32) []float32 {
	out := make([]float32, len(rgbs))
	channelLength := len(rgbs) / 3
	for i := 0; i < channelLength; i++ {
		out[i] = rgbs[i*3]
		out[i+channelLength] = rgbs[i*3+1]
		out[i+channelLength*2] = rgbs[i*3+2]
	}
	return out
}

func preprocessImage(imageFile string, inputSize int) ([]float32, float32) {
	src, err := imaging.Open(imageFile)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	var newHeight, newWidth int
	im_ratio := float32(src.Bounds().Dx()) / float32(src.Bounds().Dy())
	if im_ratio > 1 { // width > height
		newWidth = inputSize
		newHeight = int(float32(newWidth) / im_ratio)
	} else {
		newHeight = inputSize
		newWidth = int(float32(newHeight) * im_ratio)		
	}

	fmt.Println(src.Bounds(), newWidth, newHeight)

	result := imaging.Resize(src, newWidth, newHeight, imaging.Lanczos)
	//fmt.Println("resize: ", result.Rect)
	result = padBox(result)

	rgbs := make([]float32, inputSize*inputSize*3)

	j := 0
	for i := range result.Pix {
		if (i+1)%4 != 0 {
			rgbs[j] = float32(result.Pix[i])
			j++
		}
	}

	//fmt.Println(rgbs[:100])

	rgbs = TransposeRGB(rgbs)

	//fmt.Println(rgbs[:100])

	channelLength := len(rgbs) / 3
	for i := 0; i < channelLength; i++ {
		rgbs[i] = normalize(rgbs[i], 127.5, 128.0)
		rgbs[i+channelLength] = normalize(rgbs[i+channelLength], 127.5, 128.0)
		rgbs[i+channelLength*2] = normalize(rgbs[i+channelLength*2], 127.5, 128.0)
	}

	//fmt.Println("det_scale===", det_scale, float32(newHeight), float32(src.Bounds().Dy()))

	return rgbs, float32(newHeight) / float32(src.Bounds().Dy())
}

func normalize(in float32, m float32, s float32) float32 {
	return (in - m) / s
}


// 调整为方形，黑色填充
func padBox(src image.Image) *image.NRGBA {
	var maxW int

	if src.Bounds().Dx() > src.Bounds().Dy() {
		maxW = src.Bounds().Dx()
	} else {
		maxW = src.Bounds().Dy()
	}

	dst := imaging.New(maxW, maxW, color.Black)
	dst = imaging.Paste(dst, src, image.Point{0,0})

	_ = imaging.Save(dst, "/tmp/test2.jpg")

	return dst
}

// 处理推理结果
func processResult(net_outs []onnxruntime.TensorValue, det_scale float32) ([][]float32, [][]float32) {
	for i:=0;i<len(net_outs);i++ {
		fmt.Printf("Success do predict, shape : %+v, result : %+v\n", 
			net_outs[i].Shape, 
			net_outs[i].Value.([]float32)[:net_outs[i].Shape[1]], // only show one value
		)
	}

	// len(outputs)==9
	_fmc := 3
	_feat_stride_fpn := []int{8, 16, 32}
	_num_anchors := 2
	//_use_kps := true


	center_cache := make(map[string][][]float32)

	var scores_list []float32
	var bboxes_list [][]float32
	var kpss_list [][]float32

	for idx := range _feat_stride_fpn {
		stride := _feat_stride_fpn[idx]
		scores := net_outs[idx].Value.([]float32)
		bbox_preds := net_outs[idx+_fmc].Value.([]float32)
		for i := range bbox_preds { 
			bbox_preds[i] = bbox_preds[i] * float32(stride)
		}

		var kps_preds []float32 // landmark
		kps_preds = net_outs[idx+_fmc*2].Value.([]float32)
		for i := range kps_preds { 
			kps_preds[i] = kps_preds[i] * float32(stride)
		}

		height := det_model_input_size / stride
		width := det_model_input_size / stride
		key := fmt.Sprintf("%d-%d-%d", height, width, stride)
		var anchor_centers [][]float32
		if val, ok := center_cache[key]; ok {
			anchor_centers = val
		} else {
			anchor_centers = make([][]float32, height*width*_num_anchors)
			for i:=0;i<height;i++ {
				for j:=0;j<width;j++ {
					for k:=0;k<_num_anchors;k++ {
						anchor_centers[i*width*_num_anchors+j*_num_anchors+k] = []float32{float32(j*stride), float32(i*stride)}
					}
				}
			}
			//fmt.Println(stride, len(anchor_centers), anchor_centers)

			if len(center_cache)<100 {
				center_cache[key] = anchor_centers
			}		
		}

		// det_thresh == 0.5
		var pos_inds []int
		for i := range scores {
			if scores[i]>det_thresh {
				pos_inds = append(pos_inds, i)
			}
		}
		//fmt.Println(">det_thresh:", pos_inds)

		//fmt.Println("kps_preds", len(kps_preds), kps_preds[len(kps_preds)-1])

		bboxes := distance2bbox(anchor_centers, bbox_preds)
		kpss := distance2kps(anchor_centers, kps_preds)

		//fmt.Println("kpss", len(kpss), kpss[len(kpss)-1])

		for i:=range pos_inds {
			scores_list = append(scores_list, scores[pos_inds[i]])
			bboxes_list = append(bboxes_list, bboxes[pos_inds[i]])
			kpss_list = append(kpss_list, kpss[pos_inds[i]])
		}
	}

	//fmt.Println(scores_list)
	//fmt.Println("kpss_list", kpss_list)

	// 对应 detect() 后续计算

	for i := range bboxes_list {
		for j:=0;j<4;j++ {
			bboxes_list[i][j] /= det_scale
		}
		bboxes_list[i] = append(bboxes_list[i], scores_list[i])

		for j:=0;j<10;j++ {
			kpss_list[i][j] /= det_scale
		}
		kpss_list[i] = append(kpss_list[i], scores_list[i])
	}

	sort.Slice(bboxes_list, func(i, j int) bool { return bboxes_list[i][4] > bboxes_list[j][4] })
	sort.Slice(kpss_list, func(i, j int) bool { return kpss_list[i][10] > kpss_list[j][10] })

	//fmt.Println(kpss_list)

	keep := nms(bboxes_list)

	//fmt.Println(keep)

	det := make([][]float32, len(keep))
	kpss := make([][]float32, len(keep))
	for i := range keep {
		det[i] = bboxes_list[keep[i]]
		kpss[i] = kpss_list[keep[i]]
	}

	return det, kpss
}


func distance2bbox(points [][]float32, distance []float32) (ret [][]float32) {
	ret = make([][]float32, len(points))
	for i := range points {
		ret[i] = []float32{
			points[i][0] - distance[i*4+0],
			points[i][1] - distance[i*4+1],
			points[i][0] + distance[i*4+2],
			points[i][1] + distance[i*4+3],
		}
	}
	return
}

func distance2kps(points [][]float32, distance []float32) (ret [][]float32) {
	ret = make([][]float32, len(points))
	for i := range points {
		ret[i] = make([]float32, 10)
		for j:=0;j<10;j=j+2 {
			ret[i][j]   = points[i][j%2] + distance[i*10+j]
			ret[i][j+1] = points[i][j%2+1] + distance[i*10+j+1]
		} 
	}
	return
}


func max(a, b float32) float32 {
	if a>b { 
		return a
	} else {
		return b
	}
}

func min(a, b float32) float32 {
	if a<b { 
		return a
	} else {
		return b
	}
}

func nms(dets [][]float32) (ret []int) {
	if len(dets)==0 {
		return
	}

	var order []int
	areas := make([]float32, len(dets))
	for i := range dets {
		order = append(order, i)
		areas[i] = (dets[i][2] - dets[i][0] + 1) * (dets[i][3] - dets[i][1] + 1)
	}
	for len(order)>0 {
		i := order[0]
		ret = append(ret, i)

		var keep []int
		for j := range order[1:] {
			xx1 := max(dets[i][0], dets[order[j+1]][0])
			yy1 := max(dets[i][1], dets[order[j+1]][1])
			xx2 := min(dets[i][2], dets[order[j+1]][2])
			yy2 := min(dets[i][3], dets[order[j+1]][3])

			w := max(0.0, xx2 - xx1 + 1)
			h := max(0.0, yy2 - yy1 + 1)
			inter := w * h
			ovr := inter / (areas[i] + areas[order[j+1]] - inter)

			//fmt.Println(i, j, ovr)

			if ovr <= nms_thresh {
				keep = append(keep, order[j+1])
			}
		}

		order = keep
	}

	return
}

func estimate_affine() {

/*
	src := []Point2f{
		{0 , 0  },
		{10, 0  },
		{10, 10 },
		{0 , 10 },
	}

	dst := []Point2f{
		{0 , 0  },
		{10, 5  },
		{10, 10 },
		{5 , 10 },
	}
*/


	dst := []Point2f{
		{218.78867, 205.74413},
		{312.13818, 202.18082},
		{279.89087, 232.69415},
		{236.05072, 302.79538},
		{313.98624, 299.34445},
	}

	src := []Point2f{
		{38.2946, 51.6963},
		{73.5318, 51.5014},
		{56.0252, 71.7366},
		{41.5493, 92.3655},
		{70.7299, 92.2041},
	}


	pvsrc := NewPoint2fVectorFromPoints(src)
	defer pvsrc.Close()

	pvdst := NewPoint2fVectorFromPoints(dst)
	defer pvdst.Close()

	fmt.Println(pvdst.ToPoints())
	fmt.Println(pvsrc.ToPoints())

	inliers := NewMat()
	defer inliers.Close()
	method := 4 // cv2.LMEDS
	ransacProjThreshold := 3.0
	maxiters := uint(2000)
	confidence := 0.99
	refineIters := uint(10)

	m := EstimateAffinePartial2DWithParams(pvdst, pvsrc, inliers, method, ransacProjThreshold, maxiters, confidence, refineIters)
	//m := EstimateAffinePartial2D(pvdst, pvsrc)
	defer m.Close()

	printM(m)
	printM(inliers)

	fmt.Println(m.Type(), m.Step())

	v, _ := m.DataPtrFloat64()
	fmt.Println(v)	
}

func printM(m Mat) {
	for i:=0;i<m.Rows();i++ {
		for j:=0;j<m.Cols();j++ {
			fmt.Printf("%v ", m.GetDoubleAt(i, j))
		}
		fmt.Printf("\n")
	}	
}

