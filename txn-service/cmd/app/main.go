package main

import app "github.com/RubenVillalpando/stori-challenge/internal/application"

func main() {
	app := app.New()
	app.Serve()
}
