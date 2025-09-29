package usecase

import (
	"cloud_file_manager/models"
	"cloud_file_manager/repository"
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