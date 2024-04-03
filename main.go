// Code by Alexy HOUBLOUP

// Last push: 03/04/2024 - 23:30

package main

import (
	"fmt"
	"groupie-tracker/app"
)

// Start the Groupie Tracker app
func main() {
	fmt.Println("Starting Groupie Tracker app...")
	groupieApp := app.GroupieApp{}
	groupieApp.Run()
}
