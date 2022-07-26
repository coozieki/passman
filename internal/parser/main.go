package parser

import (
	bytesPackage "bytes"
	"encoding/json"
	"log"
	"passman/internal/interfaces"
)

type structForParsing struct {
	Records []interfaces.Record `json:"records"`
}

type parser struct {
}

func (p *parser) Marshal(records []interfaces.Record) []byte {
	structForParsing := structForParsing{Records: records}

	data, err := json.Marshal(structForParsing)
	if err != nil {
		log.Fatal("error while marshalling records: ", err)
	}

	return data
}

func (p *parser) Parse(bytes []byte) []interfaces.Record {
	var records structForParsing

	bytes = bytesPackage.TrimPrefix(bytes, []byte("\xef\xbb\xbf"))

	if err := json.Unmarshal(bytes, &records); err != nil {
		log.Fatal("invalid password")
	}

	return records.Records
}

func NewParser() interfaces.Parser {
	return &parser{}
}
