package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"cloud_file_manager/src/dto"
	"cloud_file_manager/src/models"
	"cloud_file_manager/src/usecase"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gin-gonic/gin"
)

type fakeUserRepo struct {
	createUserFn  func(models.User) (int, error)
	getUsersFn    func() ([]models.User, error)
	getUserByIDFn func(int) (*models.User, error)
	loginFn       func(dto.UserLoginDto) (*dto.UserResponseDto, error)
}

func (f *fakeUserRepo) CreateUser(user models.User) (int, error) {
	if f.createUserFn == nil {
		panic("unexpected CreateUser call")
	}
	return f.createUserFn(user)
}

func (f *fakeUserRepo) GetUsers() ([]models.User, error) {
	if f.getUsersFn == nil {
		panic("unexpected GetUsers call")
	}
	return f.getUsersFn()
}

func (f *fakeUserRepo) GetUserById(id int) (*models.User, error) {
	if f.getUserByIDFn == nil {
		panic("unexpected GetUserById call")
	}
	return f.getUserByIDFn(id)
}

func (f *fakeUserRepo) Login(input dto.UserLoginDto) (*dto.UserResponseDto, error) {
	if f.loginFn == nil {
		panic("unexpected Login call")
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
		panic("unexpected CreateBucket call")
	}
	return f.createBucketFn(ctx, bucket)
}

func (f *fakeAwsClient) ListBuckets(ctx context.Context) ([]types.Bucket, error) {
	if f.listBucketsFn == nil {
		panic("unexpected ListBuckets call")
	}
	return f.listBucketsFn(ctx)
}

func (f *fakeAwsClient) ListBucketItems(ctx context.Context, bucket string) ([]types.Object, error) {
	if f.listBucketItemsFn == nil {
		panic("unexpected ListBucketItems call")
	}
	return f.listBucketItemsFn(ctx, bucket)
}

func (f *fakeAwsClient) GetObject(ctx context.Context, bucket, key string, ttl int64) (*v4.PresignedHTTPRequest, error) {
	if f.getObjectFn == nil {
		panic("unexpected GetObject call")
	}
	return f.getObjectFn(ctx, bucket, key, ttl)
}

func (f *fakeAwsClient) PutObjectPresignedUrl(ctx context.Context, bucket, key string, ttl int64) (*v4.PresignedHTTPRequest, error) {
	if f.putObjectPresignedURLFn == nil {
		panic("unexpected PutObjectPresignedUrl call")
	}
	return f.putObjectPresignedURLFn(ctx, bucket, key, ttl)
}

func newUserController(repo usecase.UserRepository, aws usecase.AwsClient) UserController {
	usecaseLayer := usecase.NewUserUseCase(repo, aws)
	return NewUserController(usecaseLayer)
}

func TestUserControllerGetUsersSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeUserRepo{
		getUsersFn: func() ([]models.User, error) {
			return []models.User{{ID: 1, Name: "Ana"}}, nil
		},
	}
	awsClient := &fakeAwsClient{
		createBucketFn: func(context.Context, string) (*s3.CreateBucketOutput, error) {
			panic("aws não deve ser chamado")
		},
	}

	controller := newUserController(repo, awsClient)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/users", nil)

	controller.GetUsers(ctx)

	if recorder.Code != http.StatusOK {
		t.Fatalf("esperava status 200, veio %d", recorder.Code)
	}

	var users []models.User
	if err := json.Unmarshal(recorder.Body.Bytes(), &users); err != nil {
		t.Fatalf("não conseguiu decodificar resposta: %v", err)
	}
	if len(users) != 1 || users[0].Name != "Ana" {
		t.Fatalf("resposta inesperada: %#v", users)
	}
}

func TestUserControllerGetUsersError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeUserRepo{
		getUsersFn: func() ([]models.User, error) {
			return nil, errors.New("db error")
		},
	}
	awsClient := &fakeAwsClient{
		createBucketFn: func(context.Context, string) (*s3.CreateBucketOutput, error) {
			panic("aws não deve ser chamado")
		},
	}

	controller := newUserController(repo, awsClient)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/users", nil)

	controller.GetUsers(ctx)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("esperava status 500, veio %d", recorder.Code)
	}
}

func TestUserControllerCreateUserSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeUserRepo{
		createUserFn: func(user models.User) (int, error) {
			if user.Name != "Ana" {
				t.Fatalf("nome inesperado %s", user.Name)
			}
			return 10, nil
		},
		getUsersFn:    func() ([]models.User, error) { return nil, errors.New("unused") },
		getUserByIDFn: func(int) (*models.User, error) { return nil, errors.New("unused") },
		loginFn: func(dto.UserLoginDto) (*dto.UserResponseDto, error) {
			return nil, errors.New("unused")
		},
	}

	awsClient := &fakeAwsClient{
		createBucketFn: func(ctx context.Context, bucket string) (*s3.CreateBucketOutput, error) {
			if bucket != "myawss3bucket-90902222345-10" {
				t.Fatalf("bucket inesperado %s", bucket)
			}
			return &s3.CreateBucketOutput{}, nil
		},
		listBucketsFn:     func(context.Context) ([]types.Bucket, error) { return nil, errors.New("unused") },
		listBucketItemsFn: func(context.Context, string) ([]types.Object, error) { return nil, errors.New("unused") },
		getObjectFn: func(context.Context, string, string, int64) (*v4.PresignedHTTPRequest, error) {
			return nil, errors.New("unused")
		},
		putObjectPresignedURLFn: func(context.Context, string, string, int64) (*v4.PresignedHTTPRequest, error) {
			return nil, errors.New("unused")
		},
	}

	controller := newUserController(repo, awsClient)

	body, _ := json.Marshal(models.User{Name: "Ana", Email: "ana@example.com", Password: "pwd"})
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	ctx.Request.Header.Set("Content-Type", "application/json")

	controller.CreateUser(ctx)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("esperava status 201, veio %d", recorder.Code)
	}
}

func TestUserControllerCreateUserBadJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeUserRepo{
		getUsersFn: func() ([]models.User, error) { return nil, errors.New("unused") },
	}
	awsClient := &fakeAwsClient{
		createBucketFn: func(context.Context, string) (*s3.CreateBucketOutput, error) {
			panic("aws não deve ser chamado")
		},
	}

	controller := newUserController(repo, awsClient)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("invalid-json"))
	ctx.Request.Header.Set("Content-Type", "application/json")

	controller.CreateUser(ctx)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("esperava status 400, veio %d", recorder.Code)
	}
}

func TestUserControllerGetUserByIdNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeUserRepo{
		getUserByIDFn: func(int) (*models.User, error) {
			return nil, nil
		},
	}
	awsClient := &fakeAwsClient{
		createBucketFn: func(context.Context, string) (*s3.CreateBucketOutput, error) {
			panic("aws não deve ser chamado")
		},
	}

	controller := newUserController(repo, awsClient)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Params = gin.Params{{Key: "id", Value: "5"}}
	ctx.Request = httptest.NewRequest(http.MethodGet, "/users/5", nil)

	controller.GetUserById(ctx)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("esperava status 404, veio %d", recorder.Code)
	}

	var response map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("falha ao decodificar resposta: %v", err)
	}
	if response["Message"] != "Usuário não foi encontrado na base de dados" {
		t.Fatalf("mensagem inesperada: %v", response)
	}
}

func TestUserControllerGetUserByIdInvalidParam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeUserRepo{
		getUserByIDFn: func(int) (*models.User, error) {
			return nil, nil
		},
	}
	awsClient := &fakeAwsClient{
		createBucketFn: func(context.Context, string) (*s3.CreateBucketOutput, error) {
			panic("aws não deve ser chamado")
		},
	}

	controller := newUserController(repo, awsClient)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Params = gin.Params{{Key: "id", Value: "abc"}}
	ctx.Request = httptest.NewRequest(http.MethodGet, "/users/abc", nil)

	controller.GetUserById(ctx)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("esperava status 400, veio %d", recorder.Code)
	}
}
