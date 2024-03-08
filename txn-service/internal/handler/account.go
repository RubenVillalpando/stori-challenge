package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/RubenVillalpando/stori-challenge/internal/db"
	"github.com/RubenVillalpando/stori-challenge/internal/model"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateAccount(c *gin.Context) {
	var newAccount model.NewAccountRequest
	err := c.BindJSON(&newAccount)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userExists, err := h.db.UserExists(newAccount.Owner)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	if !userExists {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("user with id %d doesn't exist", newAccount.Owner))
	}

	id, err := h.db.CreateAccount(&newAccount)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, id)
}

func (h *Handler) DepositAmount(c *gin.Context) {
	var d model.Deposit
	err := c.BindJSON(&d)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.db.MakeDeposit(&d)
	if err == db.ErrInsuficientBalance {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("the amount entered was greater than the "))
	}
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) GenerateAndSaveReport(c *gin.Context) {
	id := c.Param("id")
	accId, err := strconv.Atoi(id)

	transactions, err := h.db.GetAccountReport(accId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	err = h.db.UploadReport(transactions, id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to upload file: %v", err))
		return
	}

	return
}
