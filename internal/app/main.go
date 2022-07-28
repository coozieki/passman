package app

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/term"
	"io"
	"log"
	"os"
	"passman/internal/encryptor"
	"passman/internal/interfaces"
	"passman/internal/parser"
	"passman/internal/providers"
	"passman/internal/providers/drive"
	"passman/internal/renderer"
	"strconv"
	"strings"
	"syscall"
)

const (
	defaultFilename = "passman_db2.txt"
)

type App interface {
	Run(args []string)
}

type app struct {
	password string

	dataProvider interfaces.DataProvider
	encryptor    interfaces.Encryptor
	parser       interfaces.Parser
	renderer     interfaces.Renderer
}

func (a *app) Run(args []string) {
	var encryptedBytes []byte

	open, err := os.Open(defaultFilename)
	if err != nil {
		encryptedBytes, err = a.dataProvider.GetFile(defaultFilename)
		if errors.Is(err, providers.ErrFileNotFound) {
			fmt.Println("Creating new database file...")
			fmt.Println("Enter file password: ")

			inputReader := bufio.NewReader(os.Stdin)
			a.password = readInput(inputReader)
			a.saveRecords([]interfaces.Record{})
			fmt.Println("Database file created")

			return
		}
		log.Fatal("error while getting file from provider: ", err)
	} else {
		encryptedBytes, err = io.ReadAll(open)
		if err != nil {
			log.Fatal("error while reading file: ", err)
		}
	}

	fmt.Print("Enter Password: ")

	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal("error while getting password: ", err)
	}

	fmt.Println()

	a.password = string(bytePassword)
	bytes := a.encryptor.Decrypt(encryptedBytes, []byte(a.password))
	records := a.parser.Parse(bytes)

	switch {
	case len(args) > 1:
		switch args[1] {
		case "create":
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

			records = append(records, *record)
			a.saveRecords(records)
			a.renderer.Render(records)

			return
		case "list":
			a.renderer.Render(records)
			return
		case "remove":
			if len(args) < 3 {
				log.Fatal("record ID is not provided for removal")
			}

			index, err := strconv.Atoi(args[2])
			if err != nil {
				log.Fatal("ID must be a numeric value")
			}

			if index < 1 || index > len(records) {
				log.Fatal("ID is out of range")
			}

			if index == len(records) {
				records = records[:len(records)-1]
			} else {
				records = append(records[:index-1], records[index:]...)
			}

			a.saveRecords(records)
			a.renderer.Render(records)

			return
		case "refresh":
			encryptedBytes, _ = a.dataProvider.GetFile(defaultFilename)
			bytes := a.encryptor.Decrypt(encryptedBytes, []byte(a.password))
			records = a.parser.Parse(bytes)
			a.renderer.Render(records)
		case "change_password":
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
			a.saveRecords(records)

			fmt.Println()
			fmt.Println("password changed successfully")

			return
		case "search":
			if len(args) < 3 {
				log.Fatal("search term is not provided")
			}

			var filteredRecords []interfaces.Record

			for _, record := range records {
				if strings.Contains(record.Name, args[2]) {
					filteredRecords = append(filteredRecords, record)
				}
			}

			a.renderer.Render(filteredRecords)

			return
		default:
			log.Fatal("unknown argument: ", args[1])
		}
	default:
		a.renderer.Render(records)
		return
	}
}

func (a *app) saveRecords(records []interfaces.Record) {
	data := a.parser.Marshal(records)
	a.dataProvider.SaveFile(defaultFilename, a.encryptor.Encrypt(data, []byte(a.password)))
}

func readInput(inputReader *bufio.Reader) string {
	input, _ := inputReader.ReadString('\n')
	return strings.TrimRight(input, "\n")
}

func NewApp() App {
	return &app{
		dataProvider: drive.NewGoogleDriveProvider(),
		encryptor:    encryptor.NewWeakEncryptor(),
		parser:       parser.NewParser(),
		renderer:     renderer.NewConsoleRenderer(),
	}
}
