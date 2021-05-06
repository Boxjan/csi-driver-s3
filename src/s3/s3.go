package s3

import (
	"github.com/minio/minio-go/v7"
	"k8s.io/klog/v2"
)

type s3Client struct {
	provider string
	endpoint string
	region   string

	clientProvider string

	S3c         *minio.Client
	ExtraClient map[string]interface{}
}

type s3Mounter struct {
}

var clientAllocMapping map[string]func(params, secrets map[string]string) (*s3Client, error)

func init() {
	clientAllocMapping = make(map[string]func(params, secrets map[string]string) (*s3Client, error))

	{
		// digitalocean
		clientAllocMapping["digitalocean"] = NewDigitaloceanS3Client
		clientAllocMapping["do"] = NewDigitaloceanS3Client

		// standard
		// use minio go api, aws-sdk is too heavy
		clientAllocMapping["s3"] = NewMinioS3Client
		clientAllocMapping[""] = NewMinioS3Client
	}
}

func NewS3Client(params, secrets map[string]string) (*s3Client, error) {
	if _, ok := clientAllocMapping[params["provider"]]; ok {
		return clientAllocMapping[params["provider"]](params, secrets)
	}

	klog.Warningf("not found provider: %s, will use standard s3 client", params["provider"])

	return clientAllocMapping["s3"](params, secrets)
}

func NewS3Mounter() {

}
