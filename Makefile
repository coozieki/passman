create-release:
	mkdir release-$(v)
	go build cmd/app/passman.go
	env GOOS=windows GOARCH=amd64 go build cmd/app/passman.go
	mv passman release-$(v)
	mv passman.exe release-$(v)
	cp credentials.json.example release-$(v)
	zip -r release-$(v).zip release-$(v)
	rm -rf release-$(v)