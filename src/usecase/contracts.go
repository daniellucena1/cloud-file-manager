package usecase

import (
	"cloud_file_manager/src/dto"
	"cloud_file_manager/src/models"
	"context"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type UserRepository interface {
	CreateUser(models.User) (int, error)
	GetUsers() ([]models.User, error)
	GetUserById(int) (*models.User, error)
	Login(dto.UserLoginDto) (*dto.UserResponseDto, error)
}

type AwsClient interface {
	CreateBucket(ctx context.Context, bucket string) (*s3.CreateBucketOutput, error)
	ListBuckets(ctx context.Context) ([]types.Bucket, error)
	ListBucketItems(ctx context.Context, bucket string) ([]types.Object, error)
	GetObject(ctx context.Context, bucket, key string, ttl int64) (*v4.PresignedHTTPRequest, error)
	PutObjectPresignedUrl(ctx context.Context, bucket, key string, ttl int64) (*v4.PresignedHTTPRequest, error)
}
