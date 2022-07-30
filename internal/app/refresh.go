package app

func (a *app) refresh() {
	encryptedBytes, _ := a.dataProvider.GetFile(defaultFilename)
	bytes := a.encryptor.Decrypt(encryptedBytes, []byte(a.password))
	a.records = a.parser.Parse(bytes)
	a.renderer.Render(a.records)
}
