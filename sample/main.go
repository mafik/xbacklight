package main

import (
	"fmt"

	"mrogalski.eu/go/xbacklight"
)

func main() {
	backlighter, err := xbacklight.NewBacklighterPrimaryScreen()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	value, err := backlighter.Get()
	if err != nil {
		fmt.Println("Couldn't query backlight:", err)
		return
	}
	fmt.Println("Current backlight:", value)
}
