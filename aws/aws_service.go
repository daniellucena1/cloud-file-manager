package aws

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

type AwsService struct {
	client *s3.Client
}

func NewAwsService(client *s3.Client) AwsService {
	return AwsService{
		client: client,
	}
}

func (as *AwsService) CreateBucket(ctx context.Context, bucketName string) (*s3.CreateBucketOutput, error) {
	output, err := as.client.CreateBucket(
		ctx,
		&s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
			CreateBucketConfiguration: &types.CreateBucketConfiguration{
				LocationConstraint: types.BucketLocationConstraint("us-east-2"),
			},
		},
	)
	if err != nil {
		var owned *types.BucketAlreadyOwnedByYou
		var exists *types.BucketAlreadyExists
		if errors.As(err, &owned) {
			log.Printf("You already own bucket %s.\n", bucketName)
			err = owned
		} else if errors.As(err, &exists) {
			log.Printf("Bucket %s already exists.\n", bucketName)
			err = exists
		}
		return nil, err
	}

	fmt.Printf("Esperando o bucket %q ser criado...\n", bucketName)

	err = s3.NewBucketExistsWaiter(as.client).Wait(
		ctx,
		&s3.HeadBucketInput{Bucket: aws.String(bucketName)},
		time.Minute,
	)
	if err != nil {
		log.Printf("Tentativa falha de esperar o bucket %s ser criado.\n", bucketName)
		return nil, err
	}

	return output, nil
}

func (as *AwsService) ListBuckets(ctx context.Context) ([]types.Bucket, error) {
	var err error
	var output *s3.ListBucketsOutput
	var buckets []types.Bucket
	bucketPaginator := s3.NewListBucketsPaginator(as.client, &s3.ListBucketsInput{})

	for bucketPaginator.HasMorePages() {
		output, err = bucketPaginator.NextPage(ctx)
		if err != nil {
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) && apiErr.ErrorCode() == "Acessop negado" {
				fmt.Println("Você não tem permissão de acessar os buckets.")
				err = apiErr
			} else {
				log.Printf("Não foi possível listar os bucket, aqui está o porque: %v\n", err)
			}
			break
		} else {
			buckets = append(buckets, output.Buckets...)
		}
	}

	return buckets, err
}
