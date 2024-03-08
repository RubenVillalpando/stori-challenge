package app

import (
	"github.com/RubenVillalpando/stori-challenge/internal/handler"
	"github.com/gin-gonic/gin"
)

type App struct {
	handler *handler.Handler
}

func New() *App {
	h := handler.New()
	return &App{
		handler: h,
	}
}

func (app *App) Serve() error {

	router := gin.Default()

	router.GET("users/:id", app.handler.GetUserById)
	router.POST("user", app.handler.CreateUser)

	router.POST("transaction", app.handler.CreateTransaction)

	router.POST("account", app.handler.CreateAccount)
	router.POST("account/deposit", app.handler.DepositAmount)
	router.POST("account/report/:id", app.handler.GenerateAndSaveReport)

	return router.Run("127.0.0.1:8080")
}
