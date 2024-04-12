package main

import "github.com/mzhn-sochi/auth-service/internal/app"

func main() {
	a, cleanup, err := app.Init()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	a.Run()
}
