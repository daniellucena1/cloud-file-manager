package usecase

import "cloud_file_manager/models"

type UserUsecase struct {

}

func NewUserUseCase() UserUsecase {
	return UserUsecase{}
}

func (uu *UserUsecase) GetUsers() ([]models.User, error) {
	return []models.User{}, nil
}