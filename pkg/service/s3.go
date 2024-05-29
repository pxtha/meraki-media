package service

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"gitlab.com/merakilab9/meradia/conf"
	"gitlab.com/merakilab9/meradia/pkg/model"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type S3Service struct {
	Client *s3.Client
	Bucket string
}

func NewS3Service(client *s3.Client, bucket string) S3ServiceInterface {
	return &S3Service{
		Client: client,
		Bucket: bucket,
	}
}

type S3PutObjectAPI interface {
	PresignPutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

type S3ServiceInterface interface {
	PreUploadMedia(key string) (string, error)
	Upload(file io.Reader, awsUrl string, contentType string, contentLength int64) (res *model.UploadDataResponse, err error)
}

func PutPresignedURL(c context.Context, api S3PutObjectAPI, input *s3.PutObjectInput) (*v4.PresignedHTTPRequest, error) {
	return api.PresignPutObject(c, input)
}

func (h S3Service) PreUploadMedia(key string) (string, error) {
	input := &s3.PutObjectInput{
		Bucket: aws.String(conf.LoadEnv().AWSBucket),
		Key:    aws.String(key),
		ACL:    types.ObjectCannedACLPublicRead,
		//ContentType: aws.String("image/jpeg"),
	}

	psClient := s3.NewPresignClient(h.Client)

	resp, err := PutPresignedURL(context.TODO(), psClient, input)
	if err != nil {
		fmt.Println("Got an error retrieving pre-signed object:")
		fmt.Println(err)
		return "", err
	}

	return resp.URL, err
}

func (h S3Service) Upload(file io.Reader, awsUrl string, contentType string, contentLength int64) (res *model.UploadDataResponse, err error) {
	req, err := http.NewRequest("PUT", awsUrl, file)
	if err != nil {
		fmt.Println("error creating request", awsUrl)
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-amz-acl", conf.LoadEnv().AWSS3ACL)
	req.ContentLength = contentLength

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("failed making request")
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", string(body))
		return nil, fmt.Errorf("failed making request to upload file, err: %v", string(body))
	}
	urlInfo, _ := url.Parse(awsUrl)
	return &model.UploadDataResponse{
		Key: urlInfo.Path,
	}, nil
}
