package usecase

import (
	"cloud_file_manager/src/dto"
	"cloud_file_manager/src/models"
	"context"
	"fmt"
	"strconv"
)

type UserUsecase struct {
	repository UserRepository
	awsService AwsClient
}

func NewUserUseCase(repo UserRepository, aws AwsClient) UserUsecase {
	return UserUsecase{
		repository: repo,
		awsService: aws,
	}
}

func (uu *UserUsecase) GetUsers() ([]models.User, error) {
	return uu.repository.GetUsers()
}

func (uu *UserUsecase) CreateUser(user models.User) (models.User, error) {

	userId, err := uu.repository.CreateUser(user)
	if err != nil {
		fmt.Println(err)
		return models.User{}, err
	}

	ctx := context.Background()
	bucketName := "myawss3bucket-90902222345-" + strconv.Itoa(userId)
	_, err = uu.awsService.CreateBucket(ctx, bucketName)
	if err != nil {
		fmt.Println(err)
		return models.User{}, err
	}

	user.ID = userId

	return user, nil
}

func (uu *UserUsecase) GetUserById(id int) (*models.User, error) {

	user, err := uu.repository.GetUserById(id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return user, nil
}

func (uu *UserUsecase) Login(userDto dto.UserLoginDto) (*dto.UserResponseDto, error) {

	user, err := uu.repository.Login(userDto)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return user, nil
}
