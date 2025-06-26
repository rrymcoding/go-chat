package main

import (
	"fmt"
)

var s server

func display_menu() {

	fmt.Println("")
	fmt.Println("")
	fmt.Println("Welcome to the menu system")
	fmt.Println("")
	fmt.Println("Select a menu option from the menu")
	fmt.Println("[0]\t Start Server")
	fmt.Println("[1]\t Send Commands to Server")
	fmt.Println("[2]\t Option No. 2")
	fmt.Println("[3]\t Option No. 3")
	fmt.Println("[4]\t Option No. 4")
	fmt.Println("[5]\t Option No. 5")
	fmt.Println("")
	fmt.Println("[M] to show menu")
	fmt.Println("[E] to exit")

}

func process_menu_selection() string {

	var option string = "-1"

	fmt.Scan(&option)

	switch option {

	case "0":
		s := newServer()
		go s.run()
		s.start()

		return option
	case "1":

		return option
	case "M":
		display_menu()
		return "M"
	case "L":
		s.getRoomsList()
	case "E":
		println("Exit Pressed")
		return option

	}

	return option
}
func display_header() {
	fmt.Print(header)

}
