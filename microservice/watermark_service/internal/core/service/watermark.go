package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math"

	"github.com/brandoyts/watmarker/microservice/watermark_service/internal/core/port"
	"github.com/fogleman/gg"
	"github.com/google/uuid"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

type WatermarkService struct {
	imageStorage port.ImageStorage
	cdnEndpoint  string
}

func NewWatermarkService(imageStorage port.ImageStorage, cdnEndpoint string) *WatermarkService {
	return &WatermarkService{
		imageStorage: imageStorage,
		cdnEndpoint:  cdnEndpoint,
	}
}

func (ws *WatermarkService) ProcessImage(ctx context.Context, in port.ProcessImageInput) (string, error) {
	// decode image data from bytes
	baseImage, baseImageExtension, err := ws.loadImage(in.ImageData)
	if err != nil {
		return "", err
	}

	rgba, err := ws.applyWatermark(baseImage, in.WatermarkText)
	if err != nil {
		return "", err
	}

	encodedImageBytes, err := ws.encodeImage(rgba, baseImageExtension)
	if err != nil {
		return "", nil
	}

	// upload image
	filename := fmt.Sprintf("%v.%v", uuid.NewString(), baseImageExtension)

	err = ws.imageStorage.Upload(ctx, filename, encodedImageBytes)
	if err != nil {
		return "", err
	}

	imageUrl := fmt.Sprintf("%v/%v", ws.cdnEndpoint, filename)

	return imageUrl, nil
}

func (ws *WatermarkService) applyWatermark(baseImage image.Image, watermarkText string) (*image.RGBA, error) {
	if baseImage == nil {
		return nil, errors.New("cannot apply watermark on nil image")
	}

	// create drawing context
	rgba := image.NewRGBA(baseImage.Bounds())
	draw.Draw(rgba, rgba.Bounds(), baseImage, image.Point{}, draw.Src)

	dc := gg.NewContextForRGBA(rgba)

	fontData, err := opentype.Parse(goregular.TTF)
	if err != nil {
		return nil, fmt.Errorf("failed to parse default font: %w", err)
	}

	imageWidth := rgba.Bounds().Dx()
	imageHeight := rgba.Bounds().Dy()
	fontSize := float64(imageHeight) * 0.20 // 20% of image height

	if fontSize < 12 {
		fontSize = 12
	}

	fontFace, err := opentype.NewFace(fontData, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create font face: %w", err)
	}

	dc.SetFontFace(fontFace)

	// configure watermark style: light gray with alpha
	r := 200.0 / 255.0
	gc := 200.0 / 255.0
	bc := 200.0 / 255.0
	alpha := 0.25 // adjust opacity
	dc.SetRGBA(r, gc, bc, alpha)

	angle := -45.0 * (math.Pi / 180.0)

	dc.Push()
	dc.RotateAbout(angle, float64(imageWidth)/2.0, float64(imageHeight)/2.0)

	dc.DrawStringAnchored(watermarkText, float64(imageWidth)/2.0, float64(imageHeight)/2.0, 0.5, 0.5)
	dc.Pop()

	return rgba, nil
}

func (ws *WatermarkService) encodeImage(rgba *image.RGBA, baseImageExtension string) ([]byte, error) {
	if rgba == nil {
		return nil, errors.New("cannot encode nil image")
	}

	var processedImageBuffer bytes.Buffer

	switch baseImageExtension {
	case "jpg", "jpeg":
		err := jpeg.Encode(&processedImageBuffer, rgba, nil)
		if err != nil {
			return nil, err
		}
	case "png":
		err := png.Encode(&processedImageBuffer, rgba)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported image format: %v", baseImageExtension)
	}

	return processedImageBuffer.Bytes(), nil
}

func (ws *WatermarkService) loadImage(imageData []byte) (image.Image, string, error) {
	imgReader := bytes.NewReader(imageData)

	baseImage, format, err := image.Decode(imgReader)
	if err != nil {
		return nil, "", err
	}

	return baseImage, format, nil
}
