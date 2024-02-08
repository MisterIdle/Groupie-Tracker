package main

import (
	"myapp/api"

	"fyne.io/fyne/v2/app"
)

func main() {
	api.CallAPI()

	a := app.New()
	w := a.NewWindow("Salut")
	w.ShowAndRun()
}
