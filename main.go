package main

import (
	"fmt"
)

func main() {

	display_header()
	display_menu()

	for {
		option := "empty"
		option = process_menu_selection()
		fmt.Println("[DEBUG OUTPUT ->M]\t", option)

		if option == "E" || option == "e" {
			break
		}
	}

}
