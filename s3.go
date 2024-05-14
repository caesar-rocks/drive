package drive

import (
	"bytes"
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3 struct {
	Key      string
	Secret   string
	Region   string
	Endpoint string
	Bucket   string
}

func (s *S3) getS3ServiceClient() (*s3.Client, error) {
	creds := credentials.NewStaticCredentialsProvider(s.Key, s.Secret, "")

	sdkConfig, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(s.Region),
		config.WithCredentialsProvider(creds),
	)
	if err != nil {
		return nil, err
	}

	svc := s3.NewFromConfig(sdkConfig, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(s.Endpoint)
		o.Region = "auto"
	})

	return svc, nil
}

func (s *S3) Put(key string, data []byte) error {
	svc, err := s.getS3ServiceClient()
	if err != nil {
		return err
	}

	_, err = svc.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})

	return err
}

func (s *S3) Get(key string) ([]byte, error) {
	svc, err := s.getS3ServiceClient()
	if err != nil {
		return nil, err
	}

	resp, err := svc.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	return buf.Bytes(), nil
}

func (s *S3) Delete(key string) error {
	svc, err := s.getS3ServiceClient()
	if err != nil {
		return err
	}

	_, err = svc.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})

	return err
}

func (s *S3) GetURL(key string) (string, error) {
	return s.Endpoint + "/" + s.Bucket + "/" + key, nil
}

func (s *S3) GetSignedURL(key string, expireSecs int64) (string, error) {
	svc, err := s.getS3ServiceClient()
	if err != nil {
		return "", err
	}
	presignClient := s3.NewPresignClient(svc)

	request, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(expireSecs * int64(time.Second))
	})

	if err != nil {
		return "", err
	}

	return request.URL, nil
}
