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

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	decoder := xml.NewDecoder(bytes.NewReader(data))

	if isStrict {
		decoder.Strict = true
	}

	err = decoder.Decode(&fb2)
	if err != nil {
		return nil, err
	}

	fb2.FixNamespaces()

	fb2.unmarshalCoverpage(data)

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

func parseImage(data []byte) string {
	result := ""
	isQuoteOpened := false
	for _, v := range data {
		if isQuoteOpened {
			if v == '"' {
				break
			}
			result += string(v)
		} else {
			if v == '"' {
				isQuoteOpened = true
			}
		}
	}
	return result
}
