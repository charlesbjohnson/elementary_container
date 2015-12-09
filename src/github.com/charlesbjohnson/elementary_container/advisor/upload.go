package advisor

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/charlesbjohnson/elementary_container/communication"
	"github.com/charlesbjohnson/elementary_container/filedigest"
	"github.com/charlesbjohnson/elementary_container/path"
	"github.com/dghubble/sling"
)

func (application *Application) upload(serverUrl string, file *os.File) (*downloader, error) {
	checksum, err := filedigest.Digest(file, sha256.New())
	if err != nil {
		return nil, err
	}

	requestBuilder := sling.New().Base(serverUrl).Path("images/")

	imageCreateInput := communication.ImageCreateInput{Checksum: checksum}
	imageCreateRequest := requestBuilder.New().Post("").BodyJSON(imageCreateInput)
	imageCreateOutput, errorMessage := new(communication.ImageCreateOutput), new(communication.Error)

	if _, err := imageCreateRequest.Receive(imageCreateOutput, errorMessage); err != nil {
		return nil, err
	}

	if errorMessage.Message != "" {
		return nil, fmt.Errorf(errorMessage.Message)
	}

	bucket := aws.String(imageCreateOutput.AWSS3Bucket)
	key := aws.String(checksum + path.FullExt(file.Name()))

	s3Service := s3.New(aws.NewConfig().WithRegion(imageCreateOutput.AWSRegion).WithCredentials(
		credentials.NewStaticCredentials(
			imageCreateOutput.AWSAccessKeyId,
			imageCreateOutput.AWSSecretAccessKey,
			imageCreateOutput.AWSSessionToken,
		)))

	if !imageCreateOutput.Exists {
		s3Uploader := s3manager.NewUploader(&s3manager.UploadOptions{S3: s3Service})

		input := &s3manager.UploadInput{Bucket: bucket, Key: key, Body: file}
		if _, err := s3Uploader.Upload(input); err != nil {
			return nil, err
		}

		application.Log.WithField("file", *key).Info("image uploaded")

		imageCommitInput := communication.ImageCommitInput{
			AWSAccessKeyId:     imageCreateOutput.AWSAccessKeyId,
			AWSSecretAccessKey: imageCreateOutput.AWSSecretAccessKey,
			AWSSessionToken:    imageCreateOutput.AWSSessionToken,
			AWSS3ObjectKey:     *key,
		}
		imageCommitRequest := requestBuilder.New().Post("commit/").BodyJSON(imageCommitInput)
		errorMessage = new(communication.Error)

		if _, err := imageCommitRequest.Receive(nil, errorMessage); err != nil {
			return nil, err
		}

		if errorMessage.Message != "" {
			return nil, fmt.Errorf(errorMessage.Message)
		}

		application.Log.WithField("file", *key).Info("image committed")
	}

	imageDownloader := &downloader{
		strategy: func(writerAt io.WriterAt) error {
			s3Downloader := s3manager.NewDownloader(&s3manager.DownloadOptions{S3: s3Service})
			input := &s3.GetObjectInput{Bucket: bucket, Key: key}

			if _, err := s3Downloader.Download(writerAt, input); err != nil {
				return err
			}

			return nil
		},
	}

	return imageDownloader, nil
}
