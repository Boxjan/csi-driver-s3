package s3

import (
	"errors"
	"github.com/digitalocean/godo"
)

var (
	ErrNoFoundDoToken = errors.New("no fount digitalocean token in secrets, it should be named `DO-API-KEY` or `DoApiKey`")
	ErrDOTokenErr     = errors.New("can not use digitalocean token to get any info")
)

const (
	ProviderDigitalocean = "digitalocean"

	ExtraConfigDigitaloceanUseCDNKey = "digitalocean-use-cdn"
)

func NewDigitaloceanS3Client(params, secrets *map[string]string) (*s3Client, error) {
	mClient, err := NewMinioS3Client(params, secrets)
	if err != nil {
		return nil, err
	}

	mClient.provider = ProviderDigitalocean

	//digitalocean token name: DO-API-KEY
	if token, ok := getStringValueFromMapTryMultipleKey(secrets, `DO-API-KEY`, `DoApiKey`); ok {
		mClient.ExtraClient["do"] = godo.NewFromToken(token)
	} else {
		return mClient, ErrNoFoundDoToken
	}

	if v, ok := getStringValueFromMapTryMultipleKey(params, `DO-USE-CDN`, `doUseCdn`); ok && v == "true" {
		mClient.ExtraConfig[ExtraConfigDigitaloceanUseCDNKey] = v
	}

	return mClient, nil
}
