package config

type Config struct {
	Local  Local
	Remote Remote
}

type Local struct {
	AWS AWS `json:"aws"`
	SSH SSH `json:"ssh"`
}

type AWS struct {
	Profile string `json:"profile"`
	Region  string `json:"region"`
	S3      S3     `json:"s3"`
}

type SSH struct {
	PathToKey string `json:"pathToKey"`
}

type S3 struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}

type Remote struct {
	ConnectionString string `json:"connectionString"`
	PortToOpen       int    `json:"portToOpen"`
}
