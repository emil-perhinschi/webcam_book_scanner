package main

import (
	"WBS/internal/opencv"
	"fmt"

	"github.com/emil-perhinschi/tap-go"
)

func main() {
	t := tap.New()
	devices, err := opencv.ListDevices()
	t.Ok(err != nil, "error is nil")
	for _, device := range devices {
		fmt.Println(">>>", device)
	}
}
