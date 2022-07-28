package main

import (
	"log"
	"os"
	"passman/internal/app"
	"path/filepath"
)

func main() {
	execFilepath, err := os.Executable()
	if err != nil {
		log.Fatal("error while getting executable path: ", err)
	}

	execPath := filepath.Dir(execFilepath)

	if err := os.Chdir(execPath); err != nil {
		log.Fatal("error while changing working dir: ", err)
	}

	app := app.NewApp()
	app.Run(os.Args)
}
