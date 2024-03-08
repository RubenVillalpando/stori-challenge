package handler

import (
	"fmt"
	"net/http"

	"github.com/RubenVillalpando/stori-challenge/internal/db"
	"github.com/RubenVillalpando/stori-challenge/internal/model"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetTransactionById(c *gin.Context) {
	c.Writer.Write([]byte("Hello World!"))
}

func (h *Handler) CreateTransaction(c *gin.Context) {
	var tr model.TransactionRequest
	err := c.BindJSON(&tr)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	exist, err := h.db.AccountsExist(tr.Origin, tr.Destination)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !exist {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("one of the accounts don't exist"))
		return
	}

	err = h.db.CreateTransaction(&tr)
	if err == db.ErrInsuficientBalance {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("insuficient balance"))
	}
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}
