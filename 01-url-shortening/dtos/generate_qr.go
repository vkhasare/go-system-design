package dtos

type GenerateQRCodeRequest struct {
	Size        int    `json:"size,omitempty"`
	ImageFormat string `json:"image_format,omitempty"`
}
