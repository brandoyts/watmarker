package request

type ApplyWatermarkRequest struct {
	WatermarkText string `json:"watermark_text"`
	FileData      []byte `json:"file_data"`
}
