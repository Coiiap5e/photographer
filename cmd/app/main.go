package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/Coiiap5e/photographer/config"
	"github.com/Coiiap5e/photographer/internal/database"
)

func main() {
	ctx := context.Background()

	dbConfig := config.LoadDBConfig()

	db, err := database.NewClient(ctx, dbConfig)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		showMenu()
		fmt.Print("Select a menu item: ")
		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			// TODO:
		case "2":
			// TODO:
		case "3":
			// TODO:
		case "4":
			// TODO:
		case "5":
			// TODO:
		case "6":
			// TODO:
		case "7":
			// TODO:
		default:
			fmt.Println("Invalid choice")
		}
		fmt.Println("")
	}
}

func showMenu() {
	fmt.Println("1. Add client")
	fmt.Println("2. Delete client")
	fmt.Println("3. Add date of shoot")
	fmt.Println("4. Delete date of shot")
	fmt.Println("5. Show list of clients")
	fmt.Println("6. Show list of dates of shoots")
	fmt.Println("7. Exit")
}
