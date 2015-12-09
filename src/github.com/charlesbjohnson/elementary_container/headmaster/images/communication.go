package images

type CreateInput struct {
	Checksum string `json:"checksum"`
}

type CreateOutput struct {
	AWSAccessKeyID     string `json:"aws_access_key_id"`
	AWSSecretAccessKey string `json:"aws_secret_key"`
	AWSSessionToken    string `json:"aws_session_token"`
	AWSRegion          string `json:"aws_region"`
	AWSS3Bucket        string `json:"aws_s3_bucket"`
	AWSS3ObjectKey     string `json:"aws_s3_object_key"`
	Exists             bool   `json:"exists"`
}
