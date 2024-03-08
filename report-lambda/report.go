package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-mail/mail"
)

type transaction struct {
	amount float64
	date   time.Time
}

type EmailSender struct {
	d *mail.Dialer
}

func NewEmailSender() *EmailSender {
	return &EmailSender{
		d: mail.NewDialer("smtp.gmail.com", 587, "rubenvillalpando1299@gmail.com", "cwom zbph ddzu zzsy"),
	}
}

func sendMail(accId string, es *EmailSender, mail string) error {

	file, err := os.Open(fmt.Sprintf("C:/Users/mrvil/projects/stori-challenge/txn-service/reports/%s.csv", accId))
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	numTransactionsDebit := 0
	numTransactionsCredit := 0
	var balance, totalDebit, totalCredit float64
	var transactions []*transaction
	txnPerMonth := make(map[time.Month]int)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		datestr := record[1]
		amount, err := strconv.ParseFloat(record[2], 2)
		date, err := time.Parse("2006-01-02 15:04:05", datestr)
		if err != nil {
			log.Fatal("Error parsing datetime:", err)
		}
		if err != nil {
			log.Fatal("couldnt process amount")
		}
		transactions = append(transactions, &transaction{amount: amount, date: date})

	}
	for _, t := range transactions {
		txnPerMonth[t.date.Month()]++
		if t.amount < 0 {
			numTransactionsDebit++
			totalDebit += t.amount
		} else {
			numTransactionsCredit++
			totalCredit += t.amount
		}
		balance += t.amount
	}
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("<p>Summary of transactions:\nTotal balance: %.2f\n</p>", balance))
	buffer.WriteString(fmt.Sprintf("<p>Average debit amount: %.2f\n</p>", totalDebit/float64(numTransactionsDebit)))
	buffer.WriteString(fmt.Sprintf("<p>Average credit amount: %.2f\n</p>", totalCredit/float64(numTransactionsCredit)))
	for month, n := range txnPerMonth {
		buffer.WriteString(fmt.Sprintf("<p>Transactions in the month of %s: %d</p>\n", month.String(), n))
	}

	err = es.sendEmail(mail, "Transaction Summary", buffer.String())
	if err != nil {
		return err
	}

	return nil
}

func (es *EmailSender) sendEmail(recipient, subject, body string) error {
	m := mail.NewMessage()
	m.SetHeader("From", "rubenvillalpando1299@gmail.com")
	m.SetHeader("To", recipient)

	m.SetHeader("Subject", subject)

	m.SetBody("text/html", body)

	return es.d.DialAndSend(m)

}

func main() {
	s := NewEmailSender()
	sendMail("1", s, "ruby1299@gmail.com")
}
