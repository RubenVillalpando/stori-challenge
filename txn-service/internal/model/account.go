package model

type Account struct {
	ID      int     `json:"id"`
	Balance float64 `json:"balance"`
	Owner   int     `json:"owner"`
}

type Deposit struct {
	Owner   int     `json:"owner"`
	Balance float64 `json:"balance"`
}

type NewAccountRequest struct {
	Owner int `json:"owner"`
}
