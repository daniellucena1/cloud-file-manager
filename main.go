package main

import (
	"cloud_file_manager/src/app"
)

func main() {
	// setup and run app
	err := app.SetupAndRunApp()
	if err != nil {
		panic(err)
	}
}