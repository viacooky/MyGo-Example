package main

import "my-app/internal/app"

func main() {
	err := app.Run()
	if err != nil {
		panic(err)
	}
}
