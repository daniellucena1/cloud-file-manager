package usecase

import (
	"cloud_file_manager/aws"
	"context"
	"fmt"
	"strconv"

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
