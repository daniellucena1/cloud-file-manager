package usecase

import (
	"cloud_file_manager/src/aws"
	"context"
	"fmt"
	"strconv"
	"strings"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type AwsUsecase struct {
	AwsService aws.AwsService
}

func NewAwsUsecase(awsService aws.AwsService) AwsUsecase {
	return AwsUsecase{
		AwsService: awsService,
	}
}

func (au *AwsUsecase) CreateBucket(userId int, name string) (*s3.CreateBucketOutput, error) {
	bucketName := name + "-" + strconv.Itoa(userId)

	ctx := context.Background()

	output, err := au.AwsService.CreateBucket(ctx, bucketName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return output, nil
}

func (au *AwsUsecase) ListBuckets() ([]types.Bucket, error) {

	ctx := context.Background()

	output, err := au.AwsService.ListBuckets(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return output, nil
}

func (au *AwsUsecase) ListBucketItems(userId int) ([]types.Object, error) {
	ctx := context.Background()

	buckets, err := au.AwsService.ListBuckets(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var suffix string
	suffix = "-" + strconv.Itoa(userId)

	var bucketName string

	for _, element := range buckets {
		if strings.HasSuffix(*element.Name, suffix) {
			bucketName = *element.Name
		}
	}

	output, err := au.AwsService.ListBucketItems(ctx, bucketName)

	return output, err
}

func (au *AwsUsecase) GetObject(userId int, objectKey string) (*v4.PresignedHTTPRequest, error) {
	ctx := context.Background()

	buckets, err := au.AwsService.ListBuckets(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var suffix string
	suffix = "-" + strconv.Itoa(userId)

	var bucketName string

	for _, element := range buckets {
		if strings.HasSuffix(*element.Name, suffix) {
			bucketName = *element.Name
		}
	}

	output, err := au.AwsService.GetObject(ctx, bucketName, objectKey, 60)

	return output, err
}

func (au *AwsUsecase) PutObject(userId int, objectKey string) (*v4.PresignedHTTPRequest, error) {
	ctx := context.Background()

	buckets, err := au.AwsService.ListBuckets(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var suffix string
	suffix = "-" + strconv.Itoa(userId)

	var bucketName string

	for _, element := range buckets {
		if strings.HasSuffix(*element.Name, suffix) {
			bucketName = *element.Name
		}
	}

	output, err := au.AwsService.PutObjectPresignedUrl(ctx, bucketName, objectKey, 60)

	return output, err
}
