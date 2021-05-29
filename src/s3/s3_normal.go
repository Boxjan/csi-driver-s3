package s3

import (
	"errors"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"net/url"
)

var (
	ErrNotFoundS3Endpoint      = errors.New("not found driver endpoint in params, it named `endpoint`")
	ErrNotFoundAccessKeyID     = errors.New("not found driver access key id in secrets, it named `ACCESS-KEY-ID` or `accessKeyId`")
	ErrNotFoundSecretAccessKey = errors.New("not found driver secret access key in secrets `SECRET-ACCESS-KEY` or `secretAccessKey`")
)

const (
	ProviderStandard    = "s3"
	ClientProviderMinio = "minio"
)

func NewMinioS3Client(params, secrets *map[string]string) (*s3Client, error) {
	accessKeyID, ok := getStringValueFromMapTryMultipleKey(secrets, "ACCESS-KEY-ID", "accessKeyId")
	if !ok {
		return nil, ErrNotFoundAccessKeyID
	}

	secretAccessKey, ok := getStringValueFromMapTryMultipleKey(secrets, "SECRET-ACCESS-KEY", "secretAccessKey")
	if !ok {
		return nil, ErrNotFoundSecretAccessKey
	}

	endpointT, ok := getStringValueFromMapTryMultipleKey(params, "ENDPOINT", "endpoint")
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

	region, _ := getStringValueFromMapTryMultipleKey(params, "REGION", "region")

	useSsl := true
	if !ok {
		useSsl = endpointU.Scheme != "http"
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

	return &s3Client{
		provider:       ProviderStandard,
		endpoint:       endpoint,
		region:         region,
		clientProvider: ClientProviderMinio,
		cl:             client,
		ExtraClient:    map[string]interface{}{},
		ExtraConfig:    map[string]string{},
	}, nil
}
