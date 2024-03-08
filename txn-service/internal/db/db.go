package db

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/RubenVillalpando/stori-challenge/internal/model"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type DB struct {
	mysql *sql.DB
	s3    *s3manager.Uploader
}

var (
	ErrInsuficientBalance = errors.New("not enough balance in the account for the transaction")
)

func New() *DB {
	// apparently not even giving full write permission to the public prevents it from throwing access denied
	// sess, err := session.NewSession(&aws.Config{
	// 	Region:      aws.String("us-east-1"),
	// 	Credentials: credentials.NewStaticCredentials("AKIAXYKJXNTB7X5EZ77L", "IOeoP7BPShXA3zqgaoQpVSGVHySVeKn4/p68yzxq", ""),
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return &DB{
		mysql: NewMySqlConnection(),
		// s3:    s3manager.NewUploader(sess),
	}
}

func (db *DB) GetUserById(id string) (*model.User, error) {
	var user model.User

	row := db.mysql.QueryRow(`
	SELECT * from Users 
	WHERE id = ?`, id)

	err := row.Scan(&user.ID, &user.Email, &user.Name)
	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (db *DB) CreateUser(u *model.NewUserRequest) (int64, error) {

	stmt, err := db.mysql.Prepare("INSERT INTO Users (name, email) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(u.Name, u.Email)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db *DB) CreateAccount(a *model.NewAccountRequest) (int64, error) {
	stmt, err := db.mysql.Prepare("INSERT INTO Accounts (owner, balance) VALUES (?, 0)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(a.Owner)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (db *DB) UserExists(id int) (bool, error) {
	var exists bool
	row := db.mysql.QueryRow("SELECT EXISTS(SELECT 1 FROM Users WHERE id = ?)", id)
	if row.Err() != nil {
		return false, row.Err()
	}

	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (db *DB) AccountsExist(origin, destination int) (bool, error) {
	var originExists, destinationExists bool

	err := db.mysql.QueryRow("SELECT EXISTS(SELECT 1 FROM Accounts WHERE id = ?) AS account_exists", origin).Scan(&originExists)
	if err != nil {
		return false, err
	}

	err = db.mysql.QueryRow("SELECT EXISTS(SELECT 1 FROM Accounts WHERE id = ?) AS account_exists", destination).Scan(&destinationExists)
	if err != nil {
		return false, err
	}

	return originExists && destinationExists, nil
}

func (db *DB) MakeDeposit(d *model.Deposit) error {
	result, err := db.mysql.Exec(`
	UPDATE Accounts
	SET balance = CASE 
		WHEN balance + ? >= 0 THEN balance + ?
		ELSE balance 
	END
	WHERE id = ?;`, d.Balance, d.Balance, d.Owner)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	if rows < 0 {
		return ErrInsuficientBalance
	}
	return nil
}

func (db *DB) CreateTransaction(t *model.TransactionRequest) error {

	var originBalance float64
	err := db.mysql.QueryRow("SELECT balance FROM Accounts WHERE id = ?", t.Origin).Scan(&originBalance)
	if err != nil {
		panic(err.Error())
	}

	if originBalance >= t.Amount {
		tx, err := db.mysql.Begin()
		if err != nil {
			panic(err.Error())
		}
		defer tx.Rollback()

		_, err = tx.Exec("UPDATE Accounts SET balance = balance - ? WHERE id = ?", t.Amount, t.Origin)
		if err != nil {
			panic(err.Error())
		}

		_, err = tx.Exec("UPDATE Accounts SET balance = balance + ? WHERE id = ?", t.Amount, t.Destination)
		if err != nil {
			panic(err.Error())
		}

		_, err = tx.Exec("INSERT INTO Transactions (amount, date, origin, destination) VALUES (?, ?, ?, ?)",
			t.Amount, time.Now(), t.Origin, t.Destination)
		if err != nil {
			panic(err.Error())
		}

		err = tx.Commit()
		if err != nil {
			panic(err.Error())
		}

		fmt.Println("Transaction completed successfully.")
	} else {
		return ErrInsuficientBalance
	}
	return nil
}

func (db *DB) GetAccountReport(accId int) ([]model.Transaction, error) {
	originstmt, err := db.mysql.Prepare(`
		SELECT t.*
		FROM Transactions t
		WHERE t.origin = (SELECT id FROM Accounts WHERE id = ?);
	`)
	if err != nil {
		log.Fatal("mysql query couldn't be prepared")
	}

	rowsOrigin, err := originstmt.Query(accId)
	if err != nil {
		return nil, err
	}
	defer rowsOrigin.Close()

	var transactions []model.Transaction
	for rowsOrigin.Next() {
		var transaction model.Transaction
		if err := rowsOrigin.Scan(&transaction.ID, &transaction.Amount, &transaction.Date, &transaction.Origin, &transaction.Destination); err != nil {
			log.Fatal("failed to scan row:", err)
		}
		transaction.Amount *= -1
		transactions = append(transactions, transaction)
	}
	if err := rowsOrigin.Err(); err != nil {
		return nil, err
	}
	destinationStmt, err := db.mysql.Prepare(`
		SELECT t.*
		FROM Transactions t
		WHERE t.destination = (SELECT id FROM Accounts WHERE id = ?);
	`)
	if err != nil {
		log.Fatal("mysql query couldn't be prepared")
	}

	rowsDest, err := destinationStmt.Query(accId)
	if err != nil {
		return nil, err
	}
	defer rowsDest.Close()
	for rowsDest.Next() {
		var transaction model.Transaction
		if err := rowsDest.Scan(&transaction.ID, &transaction.Amount, &transaction.Date, &transaction.Origin, &transaction.Destination); err != nil {
			log.Fatal("failed to scan row:", err)
		}
		transactions = append(transactions, transaction)
	}
	if err := rowsOrigin.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (db *DB) UploadReport(transactions []model.Transaction, accId string) error {
	// _, err := db.s3.Upload(&s3manager.UploadInput{
	// 	Bucket: aws.String("Reports"),
	// 	Key:    aws.String(filename),
	// 	Body:   buffer,
	// }

	file, err := os.Create(fmt.Sprintf("C:/Users/mrvil/projects/stori-challenge/txn-service/reports/%s.csv", accId))
	csvWriter := csv.NewWriter(file)
	for _, transaction := range transactions {
		if err := csvWriter.Write(transaction.ToRecord()); err != nil {
			return err
		}
	}

	csvWriter.Flush()

	return err
}
