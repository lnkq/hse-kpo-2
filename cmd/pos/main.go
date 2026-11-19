package main

import (
	"fmt"
	"os"

	"pos/internal/app"
)

func main() {
	container := app.New(os.Stdin, os.Stdout, "final_scores.csv")
	if err := container.Console.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
