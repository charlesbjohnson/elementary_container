package communication

type Error struct {
	Message string `json:"message"`
}

type ImageCreateInput struct {
	Checksum string `json:"checksum"`
}

type ImageCreateOutput struct {
	AWSAccessKeyId     string `json:"aws_access_key_id"`
	AWSSecretAccessKey string `json:"aws_secret_access_key"`
	AWSSessionToken    string `json:"aws_session_token"`
	AWSRegion          string `json:"aws_region"`
	AWSS3Bucket        string `json:"aws_s3_bucket"`
	AWSS3ObjectKey     string `json:"aws_s3_object_key"`
	Exists             bool   `json:"exists"`
}

type ImageCommitInput struct {
	AWSAccessKeyId     string `json:"aws_access_key_id"`
	AWSSecretAccessKey string `json:"aws_secret_access_key"`
	AWSSessionToken    string `json:"aws_session_token"`
	AWSS3ObjectKey     string `json:"aws_s3_object_key"`
}
