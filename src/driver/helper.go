package driver

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"strings"
)

// bucket name rule: https://docs.aws.amazon.com/AmazonS3/latest/userguide/bucketnamingrules.html
func (s *S3csi) sanitizeVolumeName(n string) string {
	name := s.bucketPrefix + n
	if len(name) < 3 || len(name) > 63 {
		hash := sha1.Sum([]byte(name))
		return hex.EncodeToString(hash[:]) // will return a 40 length string
	}
	return s.bucketPrefix + name
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
