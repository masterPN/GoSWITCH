package main

import (
	"esl-service/internal/app"
	"esl-service/internal/app/helpers"
	"runtime"
)

func main() {
	// Boost it as much as it can go ...
	runtime.GOMAXPROCS(runtime.NumCPU())

	client, err := helpers.CreateClient()
	if err != nil {
		return
	}
	defer client.Close()

	for {
		msg, err := helpers.GetMessage(client)

		if err != nil {
			break
		}

		app.Execute(client, msg)
	}
}
