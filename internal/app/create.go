package app

import (
	"bufio"
	"fmt"
	"os"
	"passman/internal/interfaces"
)

func (a *app) create() {
	record := &interfaces.Record{}
	inputReader := bufio.NewReader(os.Stdin)

	fmt.Print("Name: ")
	record.Name = readInput(inputReader)
	fmt.Print("Login: ")
	record.Login = readInput(inputReader)
	fmt.Print("Password: ")
	record.Password = readInput(inputReader)
	fmt.Print("Description: ")
	record.Description = readInput(inputReader)

	a.records = append(a.records, *record)
	a.saveRecords(a.records)
	a.renderer.Render(a.records)
}
