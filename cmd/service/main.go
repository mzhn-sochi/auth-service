package main

import "github.com/mzhn-sochi/auth-service/internal/app"

func main() {

	_, cleanup, err := app.Init()
	if err != nil {
		panic(err)
	}
	defer cleanup()

}
