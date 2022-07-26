package renderer

import (
	"fmt"
	"golang.design/x/clipboard"
	"log"
	"passman/internal/interfaces"
	"strconv"
)

type consoleRenderer struct {
}

func (r *consoleRenderer) Render(records []interfaces.Record) {
	fmt.Println("-----------------------------------")
	for i, record := range records {
		fmt.Println("ID: ", i+1)
		fmt.Println("Name: ", record.Name)
		fmt.Println("Login: ", record.Login)
		fmt.Println("Description: ", record.Description)
		fmt.Println("-----------------------------------")
	}

	err := clipboard.Init()
	if err != nil {
		log.Fatal("program failed unexpectedly")
	}

	fmt.Print("Copy password to clipboard by typing ID: ")

	var chosenID string

	_, err = fmt.Scanln(&chosenID)
	if err != nil {
		log.Fatal("ID must be numeric and existent")
	}

	chosenIDToInt, err := strconv.Atoi(chosenID)
	if err != nil || chosenIDToInt < 1 || chosenIDToInt > len(records) {
		log.Fatal("ID must be numeric and existent")
	}

	fmt.Println("Password copied to clipboard")
	changed := clipboard.Write(clipboard.FmtText, []byte(records[chosenIDToInt-1].Password))
	select {
	case <-changed:
		fmt.Println(`Password is no longer available from clipboard.`)
	}
}

func NewConsoleRenderer() interfaces.Renderer {
	return &consoleRenderer{}
}
