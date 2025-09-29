package repository

import (
	"cloud_file_manager/models"
	"database/sql"
	"fmt"
)

type UserRepository struct {
	connection *sql.DB
}

func NewUserRepository(connection *sql.DB) UserRepository {
	return UserRepository{
		connection: connection,
	}
}

func (pr *UserRepository) GetUsers() ([]models.User, error) {

	query := "SELECT id, user_name, user_email, user_password FROM user"
	rows, err := pr.connection.Query(query)
	if err != nil {
		fmt.Println(err)
		return []models.User{}, err
	}

	var userList []models.User
	var userObj models.User

	for rows.Next(){
		err = rows.Scan(
			&userObj.ID,
			&userObj.Name,
			&userObj.Email,
			&userObj.Password,
		)

		if err != nil {
			fmt.Println(err)
			return []models.User{}, err
		}

		userList = append(userList, userObj)
	}

	rows.Close()

	return userList, nil
}