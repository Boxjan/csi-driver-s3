package client

import (
	"context"
	"errors"
	"github.com/boxjan/csi-driver-s3/src/s3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"net/url"
)

// will use minio-go api, aws-sdk is too heavy for us.

var (
	ErrNotFoundS3Endpoint      = errors.New("not found driver endpoint in params, it named `endpoint`")
	ErrNotFoundAccessKeyID     = errors.New("not found driver access key id in secrets, it named `ACCESS-KEY-ID` or `accessKeyId`")
	ErrNotFoundSecretAccessKey = errors.New("not found driver secret access key in secrets `SECRET-ACCESS-KEY` or `secretAccessKey`")
)

type minioS3Client struct {
	provider string
	region   string
	S3Client *minio.Client
}

func (m *minioS3Client) GetProvider() string {
	return m.provider
}

func (m *minioS3Client) VolumeExists(name string) (bool, error) {

}

func (m *minioS3Client) CreateVolume(name string) error {
	panic("implement me")
}

func (m *minioS3Client) CleanVolume(name string) error {
	panic("implement me")
}

func (m *minioS3Client) DeleteVolume(name string) error {
	panic("implement me")
}

func (m *minioS3Client) ExtraFunc(name string, args ...interface{}) ([]interface{}, error) {
	panic("implement me")
}

func (m *minioS3Client) bucketExists(name string) (bool, error) {
	return m.S3Client.BucketExists(context.Background(), name)
}

func (m *minioS3Client) createBucket(name string) error {
	return m.S3Client.MakeBucket(context.Background(),
		name,
		minio.MakeBucketOptions{
			Region: m.region,
		})
}

func (m *minioS3Client) cleanBucket(name string) error {
	panic("implement me")
}

func (m *minioS3Client) deleteBucket(name string) error {
	panic("implement me")
}

func (m *minioS3Client) folderExistsInBucket(name string) (bool, error) {
	panic("implement me")
}

func (m *minioS3Client) createFolderInExistBucket(bucketName, name string) error {
	panic("implement me")
}

func (m *minioS3Client) cleanFolderInExistBucket(bucketName, name string) error {
	panic("implement me")
}

func (m *minioS3Client) deleteFolderInExistBucket(bucketName, name string) error {
	panic("implement me")
}

func NewMinioS3Client(params, secrets map[string]string) (S3Provider, error) {

	accessKeyID, ok := s3.getStringValueFromMapTryMultipleKey(&secrets, "ACCESS-KEY-ID", "accessKeyId")
	if !ok {
		return nil, ErrNotFoundAccessKeyID
	}

	secretAccessKey, ok := s3.getStringValueFromMapTryMultipleKey(&secrets, "SECRET-ACCESS-KEY", "secretAccessKey")
	if !ok {
		return nil, ErrNotFoundSecretAccessKey
	}

	endpointT, ok := s3.getStringValueFromMapTryMultipleKey(&params, "ENDPOINT", "endpoint")
	if !ok {
		return nil, ErrNotFoundS3Endpoint
	}

	endpointU, err := url.Parse(endpointT)
	if err != nil {
		return nil, err
	}

	endpoint := endpointU.Hostname()
	if endpointU.Port() != "" {
		endpoint = endpointU.Hostname() + ":" + endpointU.Port()
	}

	region, _ := s3.getStringValueFromMapTryMultipleKey(&params, "REGION", "region")

	useSslString, ok := s3.getStringValueFromMapTryMultipleKey(&params, "USE-SSL", "useSsl")

	var useSsl bool
	if !ok {
		useSsl = endpointU.Scheme != "http"
	} else {
		if useSslString != "false" && useSslString != "FALSE" {
			useSsl = true
		}
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:        credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure:       useSsl,
		Region:       region,
		BucketLookup: 0,
	})

	if err != nil {
		return nil, err
	}

	return &minioS3Client{
		S3Client: client,
		provider: "driver",
	}, nil
}
