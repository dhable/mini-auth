package main

import (
	"fmt"
	"net/http"

	"github.com/dhable/mini-auth/app"
)

func main() {
	app, err := app.NewApp()
	if err != nil {
		fmt.Printf("FATAL: failed to create and initialize app instance: %s", err)
		return
	}

	fmt.Println("starting listener...")
	http.ListenAndServe(":8888", app.Router)
}
