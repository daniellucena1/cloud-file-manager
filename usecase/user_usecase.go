package usecase

import (
	"cloud_file_manager/models"
	"cloud_file_manager/repository"
	"fmt"
)

type UserUsecase struct {
	repository repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository) UserUsecase {
	return UserUsecase{
		repository: repo,
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

	user.ID = userId

	return user, nil
}