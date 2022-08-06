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
	records  []interfaces.Record
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

		if err != nil {
			log.Fatal("error while getting file from provider: ", err)
		}
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
	a.records = a.parser.Parse(bytes)

	switch {
	case len(args) > 1:
		switch args[1] {
		case "create":
			a.create()
		case "list":
			a.renderer.Render(a.records)
		case "remove":
			if len(args) < 3 {
				log.Fatal("record ID is not provided for removal")
			}

			id, err := strconv.Atoi(args[2])
			if err != nil {
				log.Fatal("ID must be a numeric value")
			}

			a.remove(id)
		case "refresh":
			a.refresh()
		case "change_password":
			a.changePassword()
		case "search":
			if len(args) < 3 {
				log.Fatal("search term is not provided")
			}

			a.search(args[2])
		default:
			log.Fatal("unknown argument: ", args[1])
		}
	default:
		a.renderer.Render(a.records)
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
