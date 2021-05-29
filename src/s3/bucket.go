package s3

import (
	"bytes"
	"context"
	"fmt"
	"github.com/golang/glog"
	"github.com/minio/minio-go/v7"
	"k8s.io/klog/v2"
	"strings"
)

func (sc *s3Client) CreateVolume(ctx context.Context, volumeName string) error {
	klog.V(7).Info(ctx, volumeName)

	if exists, err := sc.VolumeExists(ctx, volumeName); err != nil {
		return err
	} else if exists {
		return nil
	}

	klog.V(5).Info(volumeName, "not exists, will create it")
	// not exists, need
	bucketName, folderName := cutVolumeName(volumeName)

	if len(folderName) != 0 {
		// in a exist bucket
		return sc.CreatePathInExistsBucket(ctx, bucketName, folderName)
	} else {
		return sc.CreateBucket(ctx, bucketName)
	}
}

func (sc *s3Client) VolumeExists(ctx context.Context, volumeName string) (bool, error) {
	klog.V(5).Infof("check volume: %s is exists")
	bucketName, folderName := cutVolumeName(volumeName)

	if len(folderName) == 0 {
		return sc.BucketExists(ctx, bucketName)
	} else {
		return sc.PathExistsInBucket(ctx, bucketName, folderName)
	}

}

func (sc *s3Client) DeleteVolume(ctx context.Context, volumeName string) error {

}

func (sc *s3Client) CreateBucket(ctx context.Context, bucketName string) error {
	return sc.cl.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
		Region:        sc.region,
		ObjectLocking: false,
	})
}

func (sc *s3Client) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	return sc.cl.BucketExists(ctx, bucketName)
}

func (sc *s3Client) CreatePathInExistsBucket(ctx context.Context, bucketName, folderName string) error {

	if exists, err := sc.PathExistsInBucket(ctx, bucketName, folderName); err != nil {
		return err
	} else if exists {
		return nil
	}

	if !strings.HasSuffix(folderName, "/") {
		folderName = folderName + "/"
	}
	_, err := sc.cl.PutObject(
		ctx, bucketName, folderName, bytes.NewReader([]byte{}), 0, minio.PutObjectOptions{},
	)

	return err
}

func (sc *s3Client) PathExistsInBucket(ctx context.Context, bucketName, path string) (bool, error) {
	sc.cl.ListObjects(ctx, bucketName, minio.ListObjectsOptions{Prefix: path})

	objInfo, err := sc.cl.StatObject(ctx, bucketName, path, minio.StatObjectOptions{})
	if err != nil {
		return false, err
	}

}

func (sc *s3Client) Clean() {

}

func (sc *s3Client) removeObjects(bucketName, prefix string) error {
	objectsCh := make(chan minio.ObjectInfo)
	var listErr error

	go func() {
		defer close(objectsCh)

		for object := range sc.cl.ListObjects(context.Background(), bucketName,
			minio.ListObjectsOptions{Prefix: prefix, Recursive: true}) {
			if object.Err != nil {
				listErr = object.Err
				return
			}
			objectsCh <- object
		}
	}()

	if listErr != nil {
		glog.Error("Error listing objects", listErr)
		return listErr
	}

	select {
	default:
		opts := minio.RemoveObjectsOptions{
			GovernanceBypass: true,
		}
		errorCh := sc.cl.RemoveObjects(context.Background(), bucketName, objectsCh, opts)
		haveErrWhenRemoveObjects := false
		for e := range errorCh {
			glog.Errorf("Failed to remove object %s, error: %s", e.ObjectName, e.Err)
			haveErrWhenRemoveObjects = true
		}
		if haveErrWhenRemoveObjects {
			return fmt.Errorf("Failed to remove all objects of bucket %s", bucketName)
		}
	}

	return nil
}

// will delete files one by one without file lock
func (sc *s3Client) removeObjectsOneByOne(bucketName, prefix string) error {
	objectsCh := make(chan minio.ObjectInfo, 1)
	removeErrCh := make(chan minio.RemoveObjectError, 1)
	var listErr error

	go func() {
		defer close(objectsCh)

		for object := range sc.cl.ListObjects(context.Background(), bucketName,
			minio.ListObjectsOptions{Prefix: prefix, Recursive: true}) {
			if object.Err != nil {
				listErr = object.Err
				return
			}
			objectsCh <- object
		}
	}()

	if listErr != nil {
		glog.Error("Error listing objects", listErr)
		return listErr
	}

	go func() {
		defer close(removeErrCh)

		for object := range objectsCh {
			err := sc.cl.RemoveObject(context.Background(), bucketName, object.Key,
				minio.RemoveObjectOptions{VersionID: object.VersionID})
			if err != nil {
				removeErrCh <- minio.RemoveObjectError{
					ObjectName: object.Key,
					VersionID:  object.VersionID,
					Err:        err,
				}
			}
		}
	}()

	haveErrWhenRemoveObjects := false
	for e := range removeErrCh {
		glog.Errorf("Failed to remove object %s, error: %s", e.ObjectName, e.Err)
		haveErrWhenRemoveObjects = true
	}
	if haveErrWhenRemoveObjects {
		return fmt.Errorf("Failed to remove all objects of path %s", bucketName)
	}

	return nil
}

//func (sc *s3Client) SetFSMeta(meta *FSMeta) error {
//	b := new(bytes.Buffer)
//	json.NewEncoder(b).Encode(meta)
//	opts := minio.PutObjectOptions{ContentType: "application/json"}
//	_, err := sc.minio.PutObject(
//		sc.ctx, meta.BucketName, path.Join(meta.Prefix, metadataName), b, int64(b.Len()), opts,
//	)
//	return err
//}
//
//func (sc *s3Client) GetFSMeta(bucketName, prefix string) (*FSMeta, error) {
//	opts := minio.GetObjectOptions{}
//	obj, err := sc.minio.GetObject(sc.ctx, bucketName, path.Join(prefix, metadataName), opts)
//	if err != nil {
//		return &FSMeta{}, err
//	}
//	objInfo, err := obj.Stat()
//	if err != nil {
//		return &FSMeta{}, err
//	}
//	b := make([]byte, objInfo.Size)
//	_, err = obj.Read(b)
//
//	if err != nil && err != io.EOF {
//		return &FSMeta{}, err
//	}
//	var meta FSMeta
//	err = json.Unmarshal(b, &meta)
//	return &meta, err
//}
