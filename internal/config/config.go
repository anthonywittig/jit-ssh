package config

type Config struct {
	Local  Local
	Remote Remote
}

type Local struct {
	S3Bucket string
	S3Key    string
}

type Remote struct {
	IP string
}
