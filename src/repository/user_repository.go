package repository

import (
	"cloud_file_manager/src/dto"
	"cloud_file_manager/src/handlers"
	"cloud_file_manager/src/models"
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	connection *sql.DB
}

func NewUserRepository(connection *sql.DB) *UserRepository {
	return &UserRepository{
		connection: connection,
	}
}

func (pr *UserRepository) GetUsers() ([]models.User, error) {

	query := "SELECT id, user_name, user_email, user_password FROM users"
	rows, err := pr.connection.Query(query)
	if err != nil {
		fmt.Println(err)
		return []models.User{}, err
	}

	var userList []models.User
	var userObj models.User

	for rows.Next() {
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

func (ur *UserRepository) CreateUser(user models.User) (int, error) {
	var id int
	query, err := ur.connection.Prepare("INSERT INTO users" +
		"(user_name, user_email, user_password)" +
		" VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	err = query.QueryRow(user.Name, user.Email, hashedPassword).Scan(&id)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	query.Close()
	return id, nil
}

func (ur *UserRepository) GetUserById(id int) (*models.User, error) {

	var user models.User

	query, err := ur.connection.Prepare("SELECT * FROM users WHERE id = $1")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = query.QueryRow(id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	query.Close()
	return &user, nil
}

func (ur *UserRepository) Login(userDto dto.UserLoginDto) (*dto.UserResponseDto, error) {
	var user models.User

	query, err := ur.connection.Prepare("SELECT id, user_name, user_email, user_password FROM users WHERE user_email = $1")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = query.QueryRow(userDto.Email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	validPassword, err := validatePassword(userDto.Password, user.Password)
	if err != nil || !validPassword {
		return nil, err
	}

	var userResponse dto.UserResponseDto
	userResponse.ID = user.ID
	userResponse.Name = user.Name
	userResponse.Email = user.Email

	token, err := handlers.CreateToken(user.Name, user.Password, user.ID)
	if err != nil {
		return nil, err
	}

	userResponse.Token = token

	query.Close()
	return &userResponse, nil
}

func validatePassword(password string, savedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(savedPassword), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}
