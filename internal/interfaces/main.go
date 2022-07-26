package interfaces

type (
	Record struct {
		Name        string
		Login       string
		Password    string
		Description string
	}

	DataProvider interface {
		GetFile(filename string) []byte
		SaveFile(filename string, data []byte)
	}
	Encryptor interface {
		Encrypt(bytes []byte, key []byte) []byte
		Decrypt(bytes []byte, key []byte) []byte
	}
	Parser interface {
		Marshal([]Record) []byte
		Parse(bytes []byte) []Record
	}
	Renderer interface {
		Render(records []Record)
	}
)
