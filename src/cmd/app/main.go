package main

import (
	"WebTic-tac-toe2/internal/di"
	"fmt"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(di.Module)
	if app != nil {
		app.Run()
	} else {
		fmt.Printf("Error with fx.New()")
	}
}
