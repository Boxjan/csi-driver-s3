package driver

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"strings"
)

const (
	VolumeIdSeparate = "#$S$#"
)

// bucket name rule: https://docs.aws.amazon.com/AmazonS3/latest/userguide/bucketnamingrules.html
func (s *S3csi) sanitizeVolumeName(volumeName string) string {
	name := s.bucketPrefix + volumeName
	if len(name) < 3 || len(name) > 63 {
		hash := sha1.Sum([]byte(name))
		name = hex.EncodeToString(hash[:]) // will return a 40 length string
	}
	return s.bucketPrefix + name
}

func (s *S3csi) encodeVolumeId(volumeName string, params *map[string]string) string {
	return s.sanitizeVolumeName(volumeName) + VolumeIdSeparate + encodeParams(params)
}

func (s *S3csi) decodeVolumeId(volumeId string) (string, *map[string]string, error) {
	sp := strings.Split(volumeId, VolumeIdSeparate)
	if len(sp) != 2 {
		return "", nil, errors.New("unknown format of volumeId")
	}
	volumeName := sp[0]
	pa, err := decodeParams(sp[1])
	return volumeName, pa, err
}

func cleanSecret(s string) string {
	if len(s) <= 3 {
		return "***"
	}

	var vv strings.Builder
	_, _ = fmt.Fprintf(&vv, "%3s", s[:3])

	for i := len(s) - 3; i > 0; i-- {
		_ = vv.WriteByte('*')
	}

	return vv.String()
}

func cleanSecretMap(map[string]string) map[string]string {
	res := make(map[string]string)
	for k, v := range res {
		res[k] = cleanSecret(v)
	}
	return res
}

func encodeParams(params *map[string]string) string {
	b, _ := json.Marshal(params)
	return base64.StdEncoding.EncodeToString(b)
}

func decodeParams(s string) (*map[string]string, error) {
	var jsonByte []byte

	if _, err := base64.StdEncoding.Decode([]byte(s), jsonByte); err != nil {
		return nil, err
	} else {
		mp := make(map[string]string)
		if err := json.Unmarshal(jsonByte, &mp); err != nil {
			return nil, err
		}
		return &mp, nil
	}
}

func CleanCreateVolumeRequestSecret(request *csi.CreateVolumeRequest) csi.CreateVolumeRequest {
	return csi.CreateVolumeRequest{
		Name:                      request.Name,
		CapacityRange:             request.CapacityRange,
		VolumeCapabilities:        request.VolumeCapabilities,
		Parameters:                request.Parameters,
		Secrets:                   cleanSecretMap(request.Secrets),
		VolumeContentSource:       request.VolumeContentSource,
		AccessibilityRequirements: request.AccessibilityRequirements,
		XXX_NoUnkeyedLiteral:      request.XXX_NoUnkeyedLiteral,
		XXX_unrecognized:          request.XXX_unrecognized,
		XXX_sizecache:             request.XXX_sizecache,
	}
}

func CleanDeleteVolumeRequestSecret(request *csi.DeleteVolumeRequest) csi.DeleteVolumeRequest {
	return csi.DeleteVolumeRequest{
		VolumeId:             request.VolumeId,
		Secrets:              cleanSecretMap(request.Secrets),
		XXX_NoUnkeyedLiteral: request.XXX_NoUnkeyedLiteral,
		XXX_unrecognized:     request.XXX_unrecognized,
		XXX_sizecache:        request.XXX_sizecache,
	}
}

func CleanCreateSnapshotRequestSecret(request *csi.CreateSnapshotRequest) csi.CreateSnapshotRequest {
	return csi.CreateSnapshotRequest{
		SourceVolumeId:       request.SourceVolumeId,
		Name:                 request.Name,
		Secrets:              cleanSecretMap(request.Secrets),
		Parameters:           request.Parameters,
		XXX_NoUnkeyedLiteral: request.XXX_NoUnkeyedLiteral,
		XXX_unrecognized:     request.XXX_unrecognized,
		XXX_sizecache:        request.XXX_sizecache,
	}
}

func CleanNodePublishVolumeRequest(request *csi.NodePublishVolumeRequest) csi.NodePublishVolumeRequest {
	return csi.NodePublishVolumeRequest{
		VolumeId:             request.VolumeId,
		PublishContext:       request.PublishContext,
		StagingTargetPath:    request.StagingTargetPath,
		TargetPath:           request.TargetPath,
		VolumeCapability:     request.VolumeCapability,
		Readonly:             request.Readonly,
		Secrets:              cleanSecretMap(request.Secrets),
		VolumeContext:        request.VolumeContext,
		XXX_NoUnkeyedLiteral: request.XXX_NoUnkeyedLiteral,
		XXX_unrecognized:     request.XXX_unrecognized,
		XXX_sizecache:        request.XXX_sizecache,
	}
}

func CleanNodeStageVolumeRequest(request *csi.NodeStageVolumeRequest) csi.NodeStageVolumeRequest {
	return csi.NodeStageVolumeRequest{
		VolumeId:             request.VolumeId,
		PublishContext:       request.PublishContext,
		StagingTargetPath:    request.StagingTargetPath,
		VolumeCapability:     request.VolumeCapability,
		Secrets:              cleanSecretMap(request.Secrets),
		VolumeContext:        request.VolumeContext,
		XXX_NoUnkeyedLiteral: request.XXX_NoUnkeyedLiteral,
		XXX_unrecognized:     request.XXX_unrecognized,
		XXX_sizecache:        request.XXX_sizecache,
	}
}

func CleanControllerExpandVolume(request *csi.ControllerExpandVolumeRequest) csi.ControllerExpandVolumeRequest {
	return csi.ControllerExpandVolumeRequest{
		VolumeId:             request.VolumeId,
		CapacityRange:        request.CapacityRange,
		Secrets:              cleanSecretMap(request.Secrets),
		VolumeCapability:     request.VolumeCapability,
		XXX_NoUnkeyedLiteral: request.XXX_NoUnkeyedLiteral,
		XXX_unrecognized:     request.XXX_unrecognized,
		XXX_sizecache:        request.XXX_sizecache,
	}
}
