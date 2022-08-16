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
)

const (
	det_model_input_size = 224
	nms_thresh = float32(0.4)
	det_thresh = float32(0.5)
)

// LD_LIBRARY_PATH=/usr/local/lib go run predict_onnx.go
func main() {
	ortEnvDet := onnxruntime.NewORTEnv(onnxruntime.ORT_LOGGING_LEVEL_VERBOSE, "development")
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

	_,_ = processResult(res, det_scale)

}


func Transpose(rgbs []float32) []float32 {
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

	fmt.Println(newWidth, newHeight)

	result := imaging.Resize(src, newWidth, newHeight, imaging.Lanczos)
	fmt.Println("resize: ", result.Rect)
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

	rgbs = Transpose(rgbs)

	//fmt.Println(rgbs[:100])

	channelLength := len(rgbs) / 3
	for i := 0; i < channelLength; i++ {
		rgbs[i] = normalize(rgbs[i], 127.5, 128.0)
		rgbs[i+channelLength] = normalize(rgbs[i+channelLength], 127.5, 128.0)
		rgbs[i+channelLength*2] = normalize(rgbs[i+channelLength*2], 127.5, 128.0)
	}

	det_scale := float32(newHeight) / float32(src.Bounds().Dy())

	fmt.Println("det_scale===", det_scale, float32(newHeight), float32(src.Bounds().Dy()))

	return rgbs, det_scale
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
func processResult(net_outs []onnxruntime.TensorValue, det_scale float32) ([]float32, error) {
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

	for idx := range _feat_stride_fpn {
		stride := _feat_stride_fpn[idx]
		scores := net_outs[idx].Value.([]float32)
		bbox_preds := net_outs[idx+_fmc].Value.([]float32)
		for i := range bbox_preds { 
			bbox_preds[i] = bbox_preds[i] * float32(stride)
		}
		//var kps_preds []float32
		//if _use_kps {
		//	kps_preds = net_outs[idx+_fmc*2].Value.([]float32)
		//	for i := range kps_preds { 
		//		kps_preds[i] = kps_preds[i] * float32(stride)
		//	}
		//}
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
		fmt.Println(">det_thresh:", pos_inds)

		bboxes := distance2bbox(anchor_centers, bbox_preds)

		//fmt.Println(bboxes[len(bboxes)-1])

		for i:=range pos_inds {
			scores_list = append(scores_list, scores[pos_inds[i]])
			bboxes_list = append(bboxes_list, bboxes[pos_inds[i]])
		}
	}

	fmt.Println(scores_list)
	fmt.Println(bboxes_list)

	// 对应 detect() 后续计算

	for i := range bboxes_list {
		bboxes_list[i][0] /= det_scale
		bboxes_list[i][1] /= det_scale
		bboxes_list[i][2] /= det_scale
		bboxes_list[i][3] /= det_scale
		bboxes_list[i] = append(bboxes_list[i], scores_list[i])
	}

	sort.Slice(bboxes_list, func(i, j int) bool { return bboxes_list[i][4] > bboxes_list[j][4] })

	fmt.Println(bboxes_list)

	keep := nms(bboxes_list)

	fmt.Println(keep)

	for i := range keep {
		fmt.Println(bboxes_list[keep[i]])
	}

	return nil, nil
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

			fmt.Println(i, j, ovr)

			if ovr <= nms_thresh {
				keep = append(keep, order[j+1])
			}
		}

		order = keep
	}

	return
}