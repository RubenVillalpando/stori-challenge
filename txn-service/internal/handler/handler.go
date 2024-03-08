package handler

import "github.com/RubenVillalpando/stori-challenge/internal/db"

type Handler struct {
	db *db.DB
}

func New() *Handler {
	return &Handler{
		db: db.New(),
	}
}
