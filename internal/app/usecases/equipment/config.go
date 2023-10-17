package equipment

type minioConfig struct {
	accessKeyID     string
	secretAccessKey string
	bucketName      string
	apiURL          string
}

func NewMinioConfig(apiURL, accessKeyID, secretAccessKey string) *minioConfig {
	return &minioConfig{
		accessKeyID:     accessKeyID,
		secretAccessKey: secretAccessKey,
		apiURL:          apiURL,
	}
}

func (cfg *minioConfig) SetBucket(bucketName string) *minioConfig {
	cfg.bucketName = bucketName
	return cfg
}
