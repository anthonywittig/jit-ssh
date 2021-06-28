package config

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
)

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

func (c Config) Hash() (string, error) {
	h := sha1.New()

	j, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("error marshalling: %s", err.Error())
	}

	if _, err := io.WriteString(h, fmt.Sprintf("json: %s", j)); err != nil {
		return "", fmt.Errorf("error writing json: %s", err.Error())
	}

	return string(h.Sum(nil)), nil
}
