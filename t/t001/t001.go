package main

import (
	"WBS/internal/device"
	"fmt"

	"github.com/emil-perhinschi/tap-go"
)

func main() {
	t := tap.New()
	devices, err := device.ListDevices()
	t.Ok(err != nil, "error is nil")
	for _, device := range devices {
		fmt.Println(">>>", device)
	}
}
