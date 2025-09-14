package s3

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestUpload_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockclient(ctrl)

	instance := &Instance{
		Bucket: "test-bucket",
		client: mockClient,
	}

	data := []byte("test content")

	mockClient.EXPECT().
		PutObject(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, input *s3.PutObjectInput, _ ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
			require.Equal(t, "test-bucket", *input.Bucket)

			buf := new(bytes.Buffer)
			_, err := buf.ReadFrom(input.Body)
			require.NoError(t, err)
			require.Equal(t, data, buf.Bytes())
			return &s3.PutObjectOutput{}, nil
		})

	err := instance.Upload(context.Background(), "", data)
	require.NoError(t, err)
}

func TestUpload_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockclient(ctrl)

	instance := &Instance{
		Bucket: "test-bucket",
		client: mockClient,
	}

	mockClient.EXPECT().
		PutObject(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("upload failed"))

	err := instance.Upload(context.Background(), "", []byte("test data"))
	require.EqualError(t, err, "upload failed")
}
