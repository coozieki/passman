create-release:
	mkdir release-$(v)
	go build cmd/app/main.go
	env GOOS=windows GOARCH=amd64 go build cmd/app/main.go
	mv main release-$(v)
	mv main.exe release-$(v)
	cp credentials.json.example release-$(v)
	zip -r release-$(v).zip release-$(v)
	rm -rf release-$(v)