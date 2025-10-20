package database

import (
	"database/sql"
	"fmt"
	"os"
	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s " + "password=%s dbname=%s sslmode=disable", 
		os.Getenv("HOST"), 
		os.Getenv("PORT"), 
		os.Getenv("DB_USER"), 
		os.Getenv("DB_PASSWORD"), 
		os.Getenv("DBNAME"),
	)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Conectado ao " + os.Getenv("DBNAME"))

	return db, nil
}