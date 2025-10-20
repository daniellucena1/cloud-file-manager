package usecase

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func TestAwsUsecaseCreateBucket(t *testing.T) {
	var captured string
	client := &fakeAwsClient{
		createBucketFn: func(ctx context.Context, bucket string) (*s3.CreateBucketOutput, error) {
			captured = bucket
			return &s3.CreateBucketOutput{}, nil
		},
	}

	usecase := NewAwsUsecase(client)

	if _, err := usecase.CreateBucket(12, "base"); err != nil {
		t.Fatalf("não esperava erro, veio %v", err)
	}

	if captured != "base-12" {
		t.Fatalf("esperava bucket base-12, veio %s", captured)
	}
}

func TestAwsUsecaseListBuckets(t *testing.T) {
	expected := []types.Bucket{
		{Name: aws.String("a")},
		{Name: aws.String("b")},
	}

	client := &fakeAwsClient{
		listBucketsFn: func(context.Context) ([]types.Bucket, error) {
			return expected, nil
		},
	}

	usecase := NewAwsUsecase(client)

	buckets, err := usecase.ListBuckets()
	if err != nil {
		t.Fatalf("não esperava erro, veio %v", err)
	}

	if !reflect.DeepEqual(buckets, expected) {
		t.Fatalf("esperava buckets %#v, veio %#v", expected, buckets)
	}
}

func TestAwsUsecaseListBucketItems(t *testing.T) {
	client := &fakeAwsClient{
		listBucketsFn: func(context.Context) ([]types.Bucket, error) {
			return []types.Bucket{
				{Name: aws.String("files-10")},
				{Name: aws.String("files-77")},
			}, nil
		},
		listBucketItemsFn: func(ctx context.Context, bucket string) ([]types.Object, error) {
			if bucket != "files-77" {
				t.Fatalf("bucket inesperado %s", bucket)
			}
			return []types.Object{{Key: aws.String("doc.txt")}}, nil
		},
	}

	usecase := NewAwsUsecase(client)

	items, err := usecase.ListBucketItems(77)
	if err != nil {
		t.Fatalf("não esperava erro, veio %v", err)
	}
	if len(items) != 1 || *items[0].Key != "doc.txt" {
		t.Fatalf("esperava um item doc.txt, veio %#v", items)
	}
}

func TestAwsUsecaseGetObject(t *testing.T) {
	client := &fakeAwsClient{
		listBucketsFn: func(context.Context) ([]types.Bucket, error) {
			return []types.Bucket{
				{Name: aws.String("files-9")},
				{Name: aws.String("files-22")},
			}, nil
		},
		getObjectFn: func(ctx context.Context, bucket, key string, ttl int64) (*v4.PresignedHTTPRequest, error) {
			if bucket != "files-22" {
				t.Fatalf("bucket inesperado %s", bucket)
			}
			if key != "photo.png" {
				t.Fatalf("key inesperada %s", key)
			}
			if ttl != 60 {
				t.Fatalf("ttl inesperado %d", ttl)
			}
			return &v4.PresignedHTTPRequest{}, nil
		},
	}

	usecase := NewAwsUsecase(client)

	if _, err := usecase.GetObject(22, "photo.png"); err != nil {
		t.Fatalf("não esperava erro, veio %v", err)
	}
}

func TestAwsUsecasePutObject(t *testing.T) {
	client := &fakeAwsClient{
		listBucketsFn: func(context.Context) ([]types.Bucket, error) {
			return []types.Bucket{
				{Name: aws.String("files-22")},
			}, nil
		},
		putObjectPresignedURLFn: func(ctx context.Context, bucket, key string, ttl int64) (*v4.PresignedHTTPRequest, error) {
			if bucket != "files-22" {
				t.Fatalf("bucket inesperado %s", bucket)
			}
			if key != "upload.bin" {
				t.Fatalf("key inesperada %s", key)
			}
			if ttl != 60 {
				t.Fatalf("ttl inesperado %d", ttl)
			}
			return &v4.PresignedHTTPRequest{}, nil
		},
	}

	usecase := NewAwsUsecase(client)

	if _, err := usecase.PutObject(22, "upload.bin"); err != nil {
		t.Fatalf("não esperava erro, veio %v", err)
	}
}
