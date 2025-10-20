package repository

import (
	"database/sql"
	"os"
	"testing"

	"cloud_file_manager/src/dto"
	"cloud_file_manager/src/models"

	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserRepositoryGetUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("não foi possível criar mock do banco: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	rows := sqlmock.NewRows([]string{"id", "user_name", "user_email", "user_password"}).
		AddRow(1, "Ana", "ana@example.com", "hash1").
		AddRow(2, "João", "joao@example.com", "hash2")

	mock.ExpectQuery("SELECT id, user_name, user_email, user_password FROM users").
		WillReturnRows(rows)

	users, err := repo.GetUsers()
	if err != nil {
		t.Fatalf("não esperava erro, veio %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("esperava 2 usuários, veio %d", len(users))
	}
	if users[0].Name != "Ana" || users[1].Email != "joao@example.com" {
		t.Fatalf("valores inesperados: %#v", users)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectativas não cumpridas: %v", err)
	}
}

func TestUserRepositoryCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("não foi possível criar mock do banco: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectPrepare("INSERT INTO users\\(user_name, user_email, user_password\\) VALUES \\(\\$1, \\$2, \\$3\\) RETURNING id").
		ExpectQuery().
		WithArgs("Ana", "ana@example.com", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(5))

	id, err := repo.CreateUser(models.User{
		Name:     "Ana",
		Email:    "ana@example.com",
		Password: "senha",
	})
	if err != nil {
		t.Fatalf("não esperava erro, veio %v", err)
	}
	if id != 5 {
		t.Fatalf("esperava id 5, veio %d", id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectativas não cumpridas: %v", err)
	}
}

func TestUserRepositoryCreateUserQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("não foi possível criar mock do banco: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectPrepare("INSERT INTO users\\(user_name, user_email, user_password\\) VALUES \\(\\$1, \\$2, \\$3\\) RETURNING id").
		ExpectQuery().
		WithArgs("Ana", "ana@example.com", sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)

	_, err = repo.CreateUser(models.User{
		Name:     "Ana",
		Email:    "ana@example.com",
		Password: "senha",
	})
	if err == nil {
		t.Fatalf("esperava erro, veio nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectativas não cumpridas: %v", err)
	}
}

func TestUserRepositoryGetUserById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("não foi possível criar mock do banco: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectPrepare("SELECT \\* FROM users WHERE id = \\$1").
		ExpectQuery().
		WithArgs(7).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_name", "user_email", "user_password"}).
			AddRow(7, "Maria", "maria@example.com", "hash"))

	user, err := repo.GetUserById(7)
	if err != nil {
		t.Fatalf("não esperava erro, veio %v", err)
	}
	if user == nil || user.Name != "Maria" {
		t.Fatalf("usuário inesperado: %#v", user)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectativas não cumpridas: %v", err)
	}
}

func TestUserRepositoryGetUserByIdNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("não foi possível criar mock do banco: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectPrepare("SELECT \\* FROM users WHERE id = \\$1").
		ExpectQuery().
		WithArgs(99).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetUserById(99)
	if err != nil {
		t.Fatalf("não esperava erro, veio %v", err)
	}
	if user != nil {
		t.Fatalf("esperava nil para usuário, veio %#v", user)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectativas não cumpridas: %v", err)
	}
}

func TestUserRepositoryLogin(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("não foi possível criar mock do banco: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	hashed, err := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("falha ao gerar hash: %v", err)
	}

	mock.ExpectPrepare("SELECT id, user_name, user_email, user_password FROM users WHERE user_email = \\$1").
		ExpectQuery().
		WithArgs("ana@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_name", "user_email", "user_password"}).
			AddRow(5, "Ana", "ana@example.com", string(hashed)))

	os.Setenv("JWT_SECRET", "secret")

	response, err := repo.Login(dto.UserLoginDto{
		Email:    "ana@example.com",
		Password: "secret",
	})
	if err != nil {
		t.Fatalf("não esperava erro, veio %v", err)
	}
	if response == nil || response.Name != "Ana" || response.Token == "" {
		t.Fatalf("resposta inesperada: %#v", response)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectativas não cumpridas: %v", err)
	}
}

func TestUserRepositoryLoginInvalidPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("não foi possível criar mock do banco: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	hashed, err := bcrypt.GenerateFromPassword([]byte("right"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("falha ao gerar hash: %v", err)
	}

	mock.ExpectPrepare("SELECT id, user_name, user_email, user_password FROM users WHERE user_email = \\$1").
		ExpectQuery().
		WithArgs("ana@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_name", "user_email", "user_password"}).
			AddRow(5, "Ana", "ana@example.com", string(hashed)))

	_, err = repo.Login(dto.UserLoginDto{
		Email:    "ana@example.com",
		Password: "wrong",
	})
	if err == nil {
		t.Fatalf("esperava erro de senha inválida, veio nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectativas não cumpridas: %v", err)
	}
}
