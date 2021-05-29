package s3

func (sc *s3Client) GetProvider() string {
	return sc.provider
}

func (sc *s3Client) GetRegion() string {
	return sc.region
}
