package usecase

import (
	"context"
	"errors"
	"testing"

	"cloud_file_manager/src/dto"
	"cloud_file_manager/src/models"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type fakeUserRepo struct {
	createUserFn  func(models.User) (int, error)
	getUsersFn    func() ([]models.User, error)
	getUserByIDFn func(int) (*models.User, error)
	loginFn       func(dto.UserLoginDto) (*dto.UserResponseDto, error)
}

func (f *fakeUserRepo) CreateUser(u models.User) (int, error) {
	if f.createUserFn == nil {
		panic("CreateUser not implemented")
	}
	return f.createUserFn(u)
}

func (f *fakeUserRepo) GetUsers() ([]models.User, error) {
	if f.getUsersFn == nil {
		panic("GetUsers not implemented")
	}
	return f.getUsersFn()
}

func (f *fakeUserRepo) GetUserById(id int) (*models.User, error) {
	if f.getUserByIDFn == nil {
		panic("GetUserById not implemented")
	}
	return f.getUserByIDFn(id)
}

func (f *fakeUserRepo) Login(input dto.UserLoginDto) (*dto.UserResponseDto, error) {
	if f.loginFn == nil {
		panic("Login not implemented")
	}
	return f.loginFn(input)
}

type fakeAwsClient struct {
	createBucketFn          func(ctx context.Context, bucket string) (*s3.CreateBucketOutput, error)
	listBucketsFn           func(ctx context.Context) ([]types.Bucket, error)
	listBucketItemsFn       func(ctx context.Context, bucket string) ([]types.Object, error)
	getObjectFn             func(ctx context.Context, bucket, key string, ttl int64) (*v4.PresignedHTTPRequest, error)
	putObjectPresignedURLFn func(ctx context.Context, bucket, key string, ttl int64) (*v4.PresignedHTTPRequest, error)
}

func (f *fakeAwsClient) CreateBucket(ctx context.Context, bucket string) (*s3.CreateBucketOutput, error) {
	if f.createBucketFn == nil {
		panic("CreateBucket not implemented")
	}
	return f.createBucketFn(ctx, bucket)
}

func (f *fakeAwsClient) ListBuckets(ctx context.Context) ([]types.Bucket, error) {
	if f.listBucketsFn == nil {
		panic("ListBuckets not implemented")
	}
	return f.listBucketsFn(ctx)
}

func (f *fakeAwsClient) ListBucketItems(ctx context.Context, bucket string) ([]types.Object, error) {
	if f.listBucketItemsFn == nil {
		panic("ListBucketItems not implemented")
	}
	return f.listBucketItemsFn(ctx, bucket)
}

func (f *fakeAwsClient) GetObject(ctx context.Context, bucket, key string, ttl int64) (*v4.PresignedHTTPRequest, error) {
	if f.getObjectFn == nil {
		panic("GetObject not implemented")
	}
	return f.getObjectFn(ctx, bucket, key, ttl)
}

func (f *fakeAwsClient) PutObjectPresignedUrl(ctx context.Context, bucket, key string, ttl int64) (*v4.PresignedHTTPRequest, error) {
	if f.putObjectPresignedURLFn == nil {
		panic("PutObjectPresignedUrl not implemented")
	}
	return f.putObjectPresignedURLFn(ctx, bucket, key, ttl)
}

func TestUserUsecaseCreateUser(t *testing.T) {
	repo := &fakeUserRepo{
		createUserFn: func(user models.User) (int, error) {
			if user.Name != "Alice" {
				t.Fatalf("esperava nome Alice, veio %q", user.Name)
			}
			return 42, nil
		},
	}

	var capturedBucket string
	awsClient := &fakeAwsClient{
		createBucketFn: func(ctx context.Context, bucket string) (*s3.CreateBucketOutput, error) {
			capturedBucket = bucket
			return &s3.CreateBucketOutput{}, nil
		},
	}

	usecase := NewUserUseCase(repo, awsClient)

	created, err := usecase.CreateUser(models.User{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "secret",
	})
	if err != nil {
		t.Fatalf("não esperava erro, veio %v", err)
	}

	if created.ID != 42 {
		t.Errorf("esperava ID 42, veio %d", created.ID)
	}

	expectedBucket := "myawss3bucket-90902222345-42"
	if capturedBucket != expectedBucket {
		t.Errorf("bucket esperado %q, veio %q", expectedBucket, capturedBucket)
	}
}

func TestUserUsecaseCreateUserRepoError(t *testing.T) {
	repoErr := errors.New("db failure")
	repo := &fakeUserRepo{
		createUserFn: func(models.User) (int, error) {
			return 0, repoErr
		},
	}

	awsClient := &fakeAwsClient{
		createBucketFn: func(ctx context.Context, bucket string) (*s3.CreateBucketOutput, error) {
			t.Fatalf("não deveria chamar CreateBucket quando o repositório falha")
			return nil, nil
		},
	}

	usecase := NewUserUseCase(repo, awsClient)

	_, err := usecase.CreateUser(models.User{})
	if !errors.Is(err, repoErr) {
		t.Fatalf("esperava erro %v, veio %v", repoErr, err)
	}
}

func TestUserUsecaseCreateUserAwsError(t *testing.T) {
	repo := &fakeUserRepo{
		createUserFn: func(models.User) (int, error) {
			return 7, nil
		},
	}

	awsErr := errors.New("aws failure")
	awsClient := &fakeAwsClient{
		createBucketFn: func(ctx context.Context, bucket string) (*s3.CreateBucketOutput, error) {
			return nil, awsErr
		},
	}

	usecase := NewUserUseCase(repo, awsClient)

	_, err := usecase.CreateUser(models.User{})
	if !errors.Is(err, awsErr) {
		t.Fatalf("esperava erro %v, veio %v", awsErr, err)
	}
}
