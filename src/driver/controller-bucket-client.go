package driver

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"
	"path"
)

func (s *S3csi) GetMeta(volumeId string) (meta map[string]string, err error) {
	meta = make(map[string]string)

	var metaRawObj *minio.Object
	metaRawObj, err = s.metaBucketClient.GetObject(context.Background(), s.meteBucketName,
		s.getMetaPath(volumeId),
		minio.GetObjectOptions{})
	if err != nil {
		return
	}

	var metaRawBytes []byte
	if metaRawBytes, err = ioutil.ReadAll(metaRawObj); err != nil {
		return
	}

	if err = json.Unmarshal(metaRawBytes, &meta); err != nil {
		return
	}

	return
}

func (s *S3csi) PutMete(volumeId string, meta map[string]string) (err error) {
	var metaRawBytes []byte
	if metaRawBytes, err = json.Marshal(meta); err != nil {
		klog.Error("marshal meta map failed: ", err)
		return err
	}

	if _, err = s.metaBucketClient.PutObject(context.Background(), s.meteBucketName,
		s.getMetaPath(volumeId),
		bytes.NewReader(metaRawBytes),
		int64(len(metaRawBytes)),
		minio.PutObjectOptions{
			SendContentMd5:   true,
			DisableMultipart: true,
			Internal:         minio.AdvancedPutOptions{},
		}); err != nil {
		klog.Error("upload to meta map to ", s.meteBucketName,
			path.Join("csi-driver-s3", "meta", volumeId), "failed: ", err)
		return
	}

	return
}

func (s *S3csi) getMetaPath(volumeId string) string {
	return path.Join("csi-driver-s3", "meta", volumeId)
}
