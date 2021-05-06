package s3

func getStringValueFromMapTryMultipleKey(mp *map[string]string, keys ...string) (string, bool) {
	for _, key := range keys {
		if v, ok := (*mp)[key]; ok {
			return v, true
		}
	}

	return "", false
}
