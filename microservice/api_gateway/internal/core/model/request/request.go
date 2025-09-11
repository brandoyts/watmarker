package request

type ApplyWatermarkRequest struct {
	Text string `json:"text"`
	Size int32  `json:"size"`
}
