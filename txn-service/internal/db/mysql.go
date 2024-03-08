package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewMySqlConnection() *sql.DB {
	connPool, err := connectToDatabase()
	if err != nil {
		log.Fatal("Failed to connect to db with err:", err)
	}
	log.Println("Connection to db successful")

	err = createTables(connPool)
	if err != nil {
		log.Fatal("Failed to create tables with err: ", err)
	}
	log.Println("Tables created successfully")

	return connPool
}

func connectToDatabase() (*sql.DB, error) {
	dataSourceName := "root:admin123@tcp(127.0.0.1:3306)/mysql"
	// dataSourceName := "admin:rupystori@tcp(stori-db.cj8ya0o8mt6p.us-east-2.rds.amazonaws.com:3306)/stori_db"
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(1 * time.Hour)

	return db, nil
}

func createTables(db *sql.DB) error {
	createUsersQuery := `
	CREATE TABLE IF NOT EXISTS Users (
		id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
		email VARCHAR(255) NOT NULL,
		name VARCHAR(255) NOT NULL,
		UNIQUE (email)
	);
    `
	createTransactionsQuery := `
	CREATE TABLE IF NOT EXISTS Transactions (
		id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
		amount DECIMAL(12, 2) NOT NULL,
		date DATETIME NOT NULL,
		origin BIGINT NOT NULL,
		destination BIGINT NOT NULL
	);
	`
	createAccountsQuery := `
	CREATE TABLE IF NOT EXISTS Accounts (
		id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
		owner BIGINT UNSIGNED NOT NULL,
		balance DECIMAL(12, 2) NOT NULL,
		CONSTRAINT accounts_owner_foreign FOREIGN KEY (owner) REFERENCES Users (id)
	);`

	_, err := db.Exec(createUsersQuery)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}
	_, err = db.Exec(createTransactionsQuery)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}
	_, err = db.Exec(createAccountsQuery)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	return nil
}
