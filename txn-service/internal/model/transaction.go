package model

import (
	"fmt"
	"strconv"
)

type Transaction struct {
	ID          int     `json:"id"`
	Amount      float64 `json:"amount"`
	Date        string  `json:"date"`
	Origin      int     `json:"origin"`
	Destination int     `json:"destination"`
}

type TransactionRequest struct {
	Amount      float64 `json:"amount"`
	Origin      int     `json:"origin"`
	Destination int     `json:"destination"`
}

func (t *Transaction) ToRecord() []string {
	return []string{
		fmt.Sprint(t.ID),
		t.Date,
		strconv.FormatFloat(t.Amount, 'f', 2, 64),
	}
}
