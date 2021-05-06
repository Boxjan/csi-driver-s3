package client

import (
	"errors"
	"fmt"
	"reflect"
)

type S3Provider interface {
	GetProvider() string

	VolumeExists(name string) (bool, error)
	CreateVolume(name string) error
	CleanVolume(name string) error
	DeleteVolume(name string) error

	ExtraFunc(name string, args ...interface{}) ([]interface{}, error)

	bucketExists(name string) (bool, error)
	createBucket(name string) error
	cleanBucket(name string) error
	deleteBucket(name string) error

	folderExistsInBucket(name string) (bool, error)
	createFolderInExistBucket(bucketName, name string) error
	cleanFolderInExistBucket(bucketName, name string) error
	deleteFolderInExistBucket(bucketName, name string) error
}

var (
	ErrNoFoundFunc     = errors.New("the func not found")
	ErrFuncArgsTooLong = errors.New("the number of params is out of index")
)

func NewS3Client(params, secrets map[string]string) (S3Provider, error) {
	if provider, ok := params["client"]; ok {
		if provider == "digitalocean" {
			return NewDigitaloceanS3Client(params, secrets)
		}
	}

	return NewMinioS3Client(params, secrets)
}

func ExtraFuncRecover() ([]interface{}, error) {
	panicReason := recover()
	return nil, errors.New(fmt.Sprint(panicReason))
}

func ExtraFunc(m interface{}, funcName string, args ...interface{}) ([]interface{}, error) {
	defer ExtraFuncRecover()

	f := reflect.ValueOf(m).MethodByName(funcName)

	in := make([]reflect.Value, len(args))
	for k, param := range args {
		in[k] = reflect.ValueOf(param)
	}

	resValue := f.Call(in)

	res := make([]interface{}, len(resValue))
	for k, r := range resValue {
		res[k] = r.Interface()
	}

	return res, nil
}
