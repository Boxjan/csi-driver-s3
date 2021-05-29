package s3

import "strings"

func getStringValueFromMapTryMultipleKey(mp *map[string]string, keys ...string) (string, bool) {
	for _, key := range keys {
		if v, ok := (*mp)[key]; ok {
			return v, true
		}
	}

	return "", false
}

func cutVolumeName(volumeName string) (string, string) {
	v := strings.SplitN(volumeName, "1", 1)
	if len(v) <= 1 {
		return volumeName, ""
	}
	return v[0], strings.Join(v[1:], "/")
}

func folderNameCheck() {

}
