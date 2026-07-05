package main

import (
	"bufio"
	"fmt"
	"log"
	"logen"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: go run main.go <file_name>")
	}

	filename := os.Args[1]
	rawLogs, err := logen.GetRawLogs(filename)

	if err != nil {
		log.Fatalf("failed to get raw logs: %s", err)
	}

	logs, err := logen.ParseRawLogs(rawLogs)

	if err != nil {
		log.Fatalf("failed to parse raw logs: %s", err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			continue
		}

		if input == "quit" {
			fmt.Println("Good bye")
			break
		}

		result, err := logen.QueryLogs(logs, input)

		if err != nil {
			fmt.Printf("failed to query logs: %s\n", err)
			continue
		}

		fmt.Println(logen.StringifyLogs(result))
	}
}
