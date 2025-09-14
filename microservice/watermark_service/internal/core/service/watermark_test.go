package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"testing"

	"github.com/brandoyts/watmarker/microservice/watermark_service/internal/core/port"
	"github.com/brandoyts/watmarker/microservice/watermark_service/internal/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func createPNGBytes(t *testing.T) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{255, 0, 0, 255})

	var buf bytes.Buffer
	require.NoError(t, png.Encode(&buf, img))

	return buf.Bytes()
}

func createJPEGBytes(t *testing.T) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{0, 255, 0, 255})

	var buf bytes.Buffer
	require.NoError(t, jpeg.Encode(&buf, img, nil))

	return buf.Bytes()
}

func TestLoadImage(t *testing.T) {

	tests := []struct {
		name        string
		imageData   []byte
		wantFormat  string
		expectError bool
	}{
		{
			name:        "success - PNG image",
			imageData:   createPNGBytes(t),
			wantFormat:  "png",
			expectError: false,
		},
		{
			name:        "success - JPEG image",
			imageData:   createJPEGBytes(t),
			wantFormat:  "jpeg",
			expectError: false,
		},
		{
			name:        "failure - corrupted image data",
			imageData:   []byte("not-an-image"),
			wantFormat:  "",
			expectError: true,
		},
		{
			name:        "failure - empty image data",
			imageData:   []byte{},
			wantFormat:  "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ws := &WatermarkService{}

			img, format, err := ws.loadImage(tt.imageData)

			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, img)
				require.Empty(t, format)
			} else {
				require.NoError(t, err)
				require.NotNil(t, img)
				require.Equal(t, tt.wantFormat, format)
			}
		})
	}
}

func TestEncodeImage(t *testing.T) {
	ws := &WatermarkService{}

	createImage := func() *image.RGBA {
		img := image.NewRGBA(image.Rect(0, 0, 1, 1))
		img.Set(0, 0, color.RGBA{255, 0, 0, 255}) // red pixel
		return img
	}

	tests := []struct {
		name      string
		rgba      *image.RGBA
		ext       string
		wantError bool
	}{
		{"Encode PNG", createImage(), "png", false},
		{"Encode JPG", createImage(), "jpg", false},
		{"Encode JPEG", createImage(), "jpeg", false},
		{"Unsupported Format", createImage(), "gif", true},
		{"Nil RGBA", nil, "png", true},
		{"Empty Image", image.NewRGBA(image.Rect(0, 0, 0, 0)), "png", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := ws.encodeImage(tt.rgba, tt.ext)

			if tt.wantError {
				require.Error(t, err, "expected error for case: %s", tt.name)
				return
			}

			require.NoError(t, err, "unexpected error for case: %s", tt.name)
			require.NotEmpty(t, encoded, "encoded bytes should not be empty for case: %s", tt.name)

			// Decode back to verify validity
			var decoded image.Image
			switch tt.ext {
			case "png":
				decoded, err = png.Decode(bytes.NewReader(encoded))
			case "jpg", "jpeg":
				decoded, err = jpeg.Decode(bytes.NewReader(encoded))
			default:
				// we skip decode for unsupported formats (already covered above)
				return
			}
			require.NoError(t, err, "decoding failed for case: %s", tt.name)

			// Bounds check (only if image is not empty)
			if tt.rgba != nil && tt.rgba.Bounds().Empty() {
				require.True(t, decoded.Bounds().Empty(), "expected empty image bounds")
			} else {
				require.Equal(t, 1, decoded.Bounds().Dx(), "expected width = 1")
				require.Equal(t, 1, decoded.Bounds().Dy(), "expected height = 1")
			}
		})
	}
}

func TestProcessImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name        string
		imageData   []byte
		setupMock   func(m *mock.MockImageStorage)
		expectedErr bool
		expectedExt string
	}{
		{
			name:      "success - valid PNG",
			imageData: createPNGBytes(t),
			setupMock: func(m *mock.MockImageStorage) {
				m.EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fileName string, data []byte) error {
						require.NotEmpty(t, data, "uploaded image must not be empty")
						require.Contains(t, fileName, ".", "file name must have extension")
						return nil
					})
			},
			expectedErr: false,
			expectedExt: ".png",
		},
		{
			name:      "failure - corrupted image bytes",
			imageData: []byte("not-an-image"),
			setupMock: func(m *mock.MockImageStorage) {
				// no upload expected, image decode should fail first
			},
			expectedErr: true,
		},
		{
			name:      "failure - upload returns error",
			imageData: createPNGBytes(t),
			setupMock: func(m *mock.MockImageStorage) {
				m.EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("upload failed"))
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockImageStorage := mock.NewMockImageStorage(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockImageStorage)
			}

			ws := NewWatermarkService(mockImageStorage, "test-endpoint")

			fileName, err := ws.ProcessImage(context.Background(), port.ProcessImageInput{
				ImageData:     tt.imageData,
				WatermarkText: "sample watermark",
			})

			if tt.expectedErr {
				require.Error(t, err)
				require.Empty(t, fileName, "filename must be empty on error")
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, fileName)
				require.Contains(t, fileName, tt.expectedExt, "fileName must preserve extension")
			}
		})
	}
}

func TestApplyWatermark(t *testing.T) {
	createBaseImage := func(width, height int) image.Image {
		img := image.NewRGBA(image.Rect(0, 0, width, height))
		draw.Draw(img, img.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)
		return img
	}

	tests := []struct {
		name          string
		width, height int
		watermarkText string
		expectChange  bool
	}{
		{"NormalImageWithText", 200, 200, "Sample Watermark", true},
		{"SmallImageWithText", 20, 20, "Tiny", true},
		{"EmptyText", 200, 200, "", false}, // expect no change
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base := createBaseImage(tt.width, tt.height)

			ws := &WatermarkService{}
			result, err := ws.applyWatermark(base, tt.watermarkText)

			require.NoError(t, err)
			require.NotNil(t, result, "resulting image should not be nil")
			require.Equal(t, tt.width, result.Bounds().Dx(), "width should remain unchanged")
			require.Equal(t, tt.height, result.Bounds().Dy(), "height should remain unchanged")
		})
	}
}
