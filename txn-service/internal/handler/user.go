package handler

import (
	"errors"
	"net/http"

	"github.com/RubenVillalpando/stori-challenge/internal/model"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetUserById(c *gin.Context) {
	id := c.Param("id")

	user, err := h.db.GetUserById(id)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) CreateUser(c *gin.Context) {
	var newUser model.NewUserRequest
	err := c.BindJSON(&newUser)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if len(newUser.Name) > 255 || len(newUser.Email) > 255 {
		c.AbortWithError(http.StatusBadRequest, errors.New("name and email must be lower than 256 characters"))
		return
	}

	id, err := h.db.CreateUser(&newUser)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, id)
}
