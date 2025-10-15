package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func InputString(prompt string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s: ", prompt)
	scanner.Scan()
	return scanner.Text()
}

func InputStringRequired(prompt string) string {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("%s: ", prompt)
		scanner.Scan()
		value := strings.TrimSpace(scanner.Text())

		if value != "" {
			return value
		}
		fmt.Println("Error: field cannot be empty")
	}
}

func InputInt(prompt string) int {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("%s: ", prompt)
		scanner.Scan()
		value, err := strconv.Atoi(scanner.Text())
		if err == nil {
			return value
		}
		fmt.Println("Error: enter a number")
	}
}

func InputId(prompt string) int {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("%s: ", prompt)
		scanner.Scan()
		value, err := strconv.Atoi(scanner.Text())
		if err == nil && value > 0 {
			return value
		}
		fmt.Println("Error: enter number > 0")
	}
}

func InputDate(prompt string) time.Time {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("%s: ", prompt)
		scanner.Scan()
		date, err := time.Parse("02.01.2006", scanner.Text())
		if err == nil {
			return date
		}
		fmt.Println("Error: use format dd.mm.yyyy (for example: 20.01.2025)")
	}
}

func InputTime(prompt string, date time.Time) time.Time {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("%s: ", prompt)
		scanner.Scan()
		timeInput := scanner.Text()
		resultTime, err := time.Parse("02.01.2006 15:04", date.Format("02.01.2006")+" "+timeInput)
		if err == nil {
			return resultTime
		}
		fmt.Println("Error: use format hh:mm (for example: 15:04)")
	}
}

func InputShootDate() (time.Time, time.Time, time.Time) {
	date := InputDate("Shoot date")
	startDate := InputTime("Start time of date", date)
	endDate := InputTime("End time of date", date)

	return date, startDate, endDate
}
