package s3configurer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/anthonywittig/jit-ssh/internal/config"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

type S3Configurer struct {
	s3          *s3.Client
	localConfig config.Local
}

func New(ctx context.Context, localConfig []byte) (*S3Configurer, error) {
	s3, err := s3Client(ctx)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting s3 client: %s", err.Error()))
	}

	lc, err := parseLocalConfig(localConfig)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error parsing local config: %s", err.Error()))
	}

	return &S3Configurer{
		s3:          s3,
		localConfig: lc,
	}, nil
}

func (c *S3Configurer) GetConfig() (config.Config, error) {
	remote, err := c.getRemoteConfig()
	if err != nil {
		return config.Config{}, errors.New(fmt.Sprintf("error getting remote config: %s", err.Error()))
	}

	return config.Config{
		Local:  c.localConfig,
		Remote: remote,
	}, nil
}

func (c *S3Configurer) getRemoteConfig(ctx context.Context) (config.Remote, error) {
	out, err := c.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.localConfig.S3Bucket),
		Key:    aws.String(c.localConfig.S3Key),
	})
	if err != nil {
		return config.Remote{}, errors.New(fmt.Sprintf("error getting s3 object: %s", err.Error()))
	}

	rc := config.Remote{}
	if err := json.NewDecoder(out.Body).Decode(&rc); err != nil {
		return config.Remote{}, errors.New(fmt.Sprintf("error decoding s3 object: %s", err.Error()))
	}

	return rc, nil
}

func parseLocalConfig(localConfig []byte) (config.Local, error) {
	lc := config.Local{}
	if err := json.Unmarshal(localConfig, &lc); err != nil {
		return config.Local{}, errors.New(fmt.Sprintf("error unmarshalling local config: %s", err.Error()))
	}
	return lc, nil
}

func s3Client(ctx context.Context) (*s3.Client, error) {
	cfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting aws config: %s", err.Error()))
	}

	client := s3.NewFromConfig(cfg)

	return client, nil
}
