// Code by Alexy HOUBLOUP

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
