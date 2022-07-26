package renderer

import (
	"fmt"
	"passman/internal/interfaces"
)

type consoleRenderer struct {
}

func (r *consoleRenderer) Render(records []interfaces.Record) {
	fmt.Println("-----------------------------------")
	for i, record := range records {
		fmt.Println("ID: ", i+1)
		fmt.Println("Name: ", record.Name)
		fmt.Println("Login: ", record.Login)
		fmt.Println("Password: ", record.Password)
		fmt.Println("Description: ", record.Description)
		fmt.Println("-----------------------------------")
	}
}

func NewConsoleRenderer() interfaces.Renderer {
	return &consoleRenderer{}
}
