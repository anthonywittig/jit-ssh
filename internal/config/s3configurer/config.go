package s3configurer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/anthonywittig/jit-ssh/internal/config"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

type S3Configurer struct {
	s3 *s3.Client
}

func New(ctx context.Context) (*S3Configurer, error) {
	lc, err := getLocalConfig()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error parsing local config: %s", err.Error()))
	}

	s3, err := s3Client(ctx, lc)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting s3 client: %s", err.Error()))
	}

	return &S3Configurer{
		s3: s3,
	}, nil
}

func (c *S3Configurer) GetConfig(ctx context.Context) (config.Config, error) {
	local, err := getLocalConfig()
	if err != nil {
		return config.Config{}, errors.New(fmt.Sprintf("error getting local config: %s", err.Error()))
	}

	remote, err := c.getRemoteConfig(ctx, local)
	if err != nil {
		return config.Config{}, errors.New(fmt.Sprintf("error getting remote config: %s", err.Error()))
	}

	conf := config.Config{
		Local:  local,
		Remote: remote,
	}
	log.Printf("config is: %+v", conf)
	return conf, nil
}

func (c *S3Configurer) getRemoteConfig(ctx context.Context, local config.Local) (config.Remote, error) {
	out, err := c.s3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(local.AWS.S3.Bucket),
		Key:    aws.String(local.AWS.S3.Key),
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

func getLocalConfig() (config.Local, error) {
	localConfig, err := ioutil.ReadFile(".env.json")
	if err != nil {
		return config.Local{}, errors.New(fmt.Sprintf("error reading file: %s", err.Error()))
	}

	lc := config.Local{}
	if err := json.Unmarshal(localConfig, &lc); err != nil {
		return config.Local{}, errors.New(fmt.Sprintf("error unmarshalling local config: %s", err.Error()))
	}
	return lc, nil
}

func s3Client(ctx context.Context, lc config.Local) (*s3.Client, error) {
	region := awsconfig.WithRegion(lc.AWS.Region)
	cfg, err := awsconfig.LoadDefaultConfig(ctx, region)
	if lc.AWS.Profile != "" {
		cfg, err = awsconfig.LoadDefaultConfig(ctx, region, awsconfig.WithSharedConfigProfile(lc.AWS.Profile))
	}
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error getting aws config: %s", err.Error()))
	}

	client := s3.NewFromConfig(cfg)

	return client, nil
}
