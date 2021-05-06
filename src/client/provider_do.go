package client

import (
	"errors"
	"github.com/digitalocean/godo"
)

var (
	ErrNoFoundDoToken = errors.New("no fount digitalocean token in secrets, it should be named `DO-API-KEY` or DoApiKey")
	ErrDOTokenErr     = errors.New("can not use digitalocean token to get any info")
)

type doS3Client struct {
	provider    string
	n           *minioS3Client
	doApiClient *godo.Client
}

func (d *doS3Client) GetProvider() string {
	return d.provider
}

func (d *doS3Client) VolumeExists(name string) (bool, error) {
	return d.n.VolumeExists(name)
}

func (d *doS3Client) CreateVolume(name string) error {
	return d.n.CreateVolume(name)
}

func (d *doS3Client) CleanVolume(name string) error {
	return d.n.CleanVolume(name)
}

func (d *doS3Client) DeleteVolume(name string) error {
	return d.n.DeleteVolume(name)
}

func (d *doS3Client) ExtraFunc(name string, args ...interface{}) ([]interface{}, error) {

	return ExtraFunc(d, name, args...)
}

func (d *doS3Client) bucketExists(name string) (bool, error) {
	return d.n.bucketExists(name)
}

func (d *doS3Client) createBucket(name string) error {
	return d.n.createBucket(name)
}

func (d *doS3Client) cleanBucket(name string) error {
	return d.n.cleanBucket(name)
}

func (d *doS3Client) deleteBucket(name string) error {
	return d.n.deleteBucket(name)
}

func (d *doS3Client) folderExistsInBucket(name string) (bool, error) {
	return d.n.folderExistsInBucket(name)
}

func (d *doS3Client) createFolderInExistBucket(bucketName, name string) error {
	return d.n.createFolderInExistBucket(bucketName, name)
}

func (d *doS3Client) cleanFolderInExistBucket(bucketName, name string) error {
	return d.n.cleanFolderInExistBucket(bucketName, name)
}

func (d *doS3Client) deleteFolderInExistBucket(bucketName, name string) error {
	return d.n.deleteFolderInExistBucket(bucketName, name)
}

func NewDigitaloceanS3Client(params, secrets map[string]string) (S3Provider, error) {

	tempClient, err := NewMinioS3Client(params, secrets)
	if err != nil {
		return nil, err
	}

	mC := tempClient.(*minioS3Client)

	var doClient *godo.Client

	//digitalocean token name: DO-API-KEY
	if token, ok := secrets["DO-API-KEY"]; ok {
		doClient = godo.NewFromToken(token)

	} else if token, ok := secrets["DoApiKey"]; ok {
		doClient = godo.NewFromToken(token)

	} else {
		return nil, ErrNoFoundDoToken
	}

	// Should we need verify token here?

	return &doS3Client{
		provider:    "do",
		n:           mC,
		doApiClient: doClient,
	}, nil
}
