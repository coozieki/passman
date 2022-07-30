package app

import (
	"fmt"
	"golang.org/x/term"
	"log"
	"syscall"
)

func (a *app) changePassword() {
	fmt.Print("New password: ")
	newPass, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal("error while getting new password: ", err)
	}

	fmt.Println()
	fmt.Print("Confirm password: ")
	confirmPass, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal("error while getting confirm password: ", err)
	}

	if string(newPass) != string(confirmPass) {
		log.Fatal("passwords do not match")
	}

	a.password = string(newPass)
	a.saveRecords(a.records)

	fmt.Println()
	fmt.Println("password changed successfully")
}
