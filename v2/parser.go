package gofb2

import (
	"bytes"
	"encoding/xml"
	"io"
	"os"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Unmarshal(reader io.Reader, isStrict bool) (*FB2, error) {
	var fb2 FB2

	decoder := xml.NewDecoder(reader)

	if isStrict {
		decoder.Strict = true
	}

	err := decoder.Decode(&fb2)
	if err != nil {
		return nil, err
	}

	fb2.FixNamespaces()

	return &fb2, nil
}

func (p *Parser) UnmarshalFromFile(filePath string, isStrict bool) (*FB2, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return p.Unmarshal(file, isStrict)
}

func (p *Parser) Marshal(fb2 *FB2, pretty bool) ([]byte, error) {
	var buf bytes.Buffer

	encoder := xml.NewEncoder(&buf)
	if pretty {
		encoder.Indent("", "  ")
	}

	fb2.FixNamespaces()

	err := encoder.Encode(fb2)
	if err != nil {
		return nil, err
	}

	err = encoder.Flush()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (p *Parser) MarshalToFile(fb2 *FB2, filePath string, pretty bool) error {
	xmlData, err := p.Marshal(fb2, pretty)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, xmlData, 0644)
}
