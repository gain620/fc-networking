package main

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
)

func PrintLogo() {
	//dev := figure.NewColorFigure("Gain", "", "green", true)
	//dev.Print()
	app := figure.NewColorFigure("FC-VM", "", "blue", true)
	app.Print()
	ver := figure.NewColorFigure("1.0.0", "", "red", true)
	ver.Print()
	fmt.Println()
}
