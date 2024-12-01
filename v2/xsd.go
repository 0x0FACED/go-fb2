package gofb2

import (
	"encoding/xml"
	"errors"
)

// Пространства имен
const (
	XMLNS_FICTIONBOOK = "http://www.gribuser.ru/xml/fictionbook/2.0"
	XMLNS_XLINK       = "http://www.w3.org/1999/xlink"
	XMLNS_GENRE       = "http://www.gribuser.ru/xml/fictionbook/2.0/genres"
	XMLNS_XS          = "http://www.w3.org/2001/XMLSchema"
)

// FB2 описывает корневую структуру FB2 файла
type FB2 struct {
	XMLName     xml.Name    `xml:"FictionBook"`
	XMLNS       *string     `xml:"xmlns,attr"`       // Основное пространство имен
	XMLNSXLink  *string     `xml:"xmlns:l,attr"`     // Пространство имен для ссылок
	XMLNSGenre  *string     `xml:"xmlns:genre,attr"` // Пространство имен для жанров
	XMLNSXS     *string     `xml:"xmlns:xs,attr"`    // Пространство имен для схемы
	Stylesheet  []string    `xml:"stylesheet"`
	Description Description `xml:"description"`
	Bodies      []Body      `xml:"body"`
	Binary      []Binary    `xml:"binary"`
}

// Body содержит секции текста книги
type Body struct {
	Sections []Section `xml:"section"`
	Name     string    `xml:"name,attr,omitempty"`
}

// Section представляет собой раздел текста книги
type Section struct {
	Title       Title     `xml:"title,omitempty"`
	P           []string  `xml:"p"`
	Subsections []Section `xml:"section,omitempty"` // Добавленное для вложенных секций
}

type Paragraph struct {
	Text string `xml:",chardata"`
}

type Title struct {
	Paragraphs []Paragraph `xml:"p,omitempty"`
}

// Binary описывает бинарные данные (например, изображения)
type Binary struct {
	Value       string `xml:",chardata"`
	ContentType string `xml:"content-type,attr"`
	ID          string `xml:"id,attr"`
}

// Description содержит метаданные книги
type Description struct {
	TitleInfo    TitleInfo    `xml:"title-info"`
	Coverpage    Coverpage    `xml:"coverpage"`
	DocumentInfo DocumentInfo `xml:"document-info"`
	PublishInfo  PublishInfo  `xml:"publish-info,omitempty"`
	CustomInfo   []CustomInfo `xml:"custom-info,omitempty"`
}

// TitleInfo содержит основную информацию о книге
type TitleInfo struct {
	Genre      []string     `xml:"genre"`
	Author     []AuthorType `xml:"author"`
	BookTitle  string       `xml:"book-title"`
	Annotation string       `xml:"annotation,omitempty"`
	Keywords   string       `xml:"keywords,omitempty"`
	Date       Date         `xml:"date"`
	Lang       string       `xml:"lang"`
	SrcLang    string       `xml:"src-lang,omitempty"`
	Translator []AuthorType `xml:"translator,omitempty"`
	Sequence   Sequence     `xml:"sequence,omitempty"`
}

// AuthorType описывает автора книги
type AuthorType struct {
	FirstName  string `xml:"first-name"`
	MiddleName string `xml:"middle-name,omitempty"`
	LastName   string `xml:"last-name"`
	Nickname   string `xml:"nickname,omitempty"`
	HomePage   string `xml:"home-page,omitempty"`
	Email      string `xml:"email,omitempty"`
}

// Coverpage описывает обложку книги
type Coverpage struct {
	Image *Image `xml:"image"`
}

// Image представляет изображение с поддержкой пространства имен xlink
type Image struct {
	Href string `xml:"xlink:href,attr"`
}

// DocumentInfo содержит информацию о документе
type DocumentInfo struct {
	Author      []AuthorType `xml:"author"`
	ProgramUsed string       `xml:"program-used,omitempty"`
	Date        Date         `xml:"date"`
	SrcURL      []string     `xml:"src-url,omitempty"`
	SrcOcr      string       `xml:"src-ocr,omitempty"`
	ID          string       `xml:"id"`
	Version     float64      `xml:"version"`
	History     string       `xml:"history,omitempty"`
}

// PublishInfo содержит издательскую информацию
type PublishInfo struct {
	BookName  string   `xml:"book-name"`
	Publisher string   `xml:"publisher"`
	City      string   `xml:"city"`
	Year      int      `xml:"year,omitempty"`
	ISBN      string   `xml:"isbn,omitempty"`
	Sequence  Sequence `xml:"sequence,omitempty"`
}

// CustomInfo содержит пользовательскую информацию
type CustomInfo struct {
	InfoType string `xml:"info-type,attr"`
	Value    string `xml:",chardata"`
}

// Sequence описывает серию книги
type Sequence struct {
	Name   string `xml:"name,attr"`
	Number int    `xml:"number,attr,omitempty"`
}

// Date описывает дату с атрибутами
type Date struct {
	Value string `xml:"value,attr,omitempty"`
	Text  string `xml:",chardata"`
}

func (fb2 FB2) Sections() []*Section {
	var sections []*Section

	for _, body := range fb2.Bodies {
		for _, section := range body.Sections {
			sections = append(sections, &section)
		}
	}

	return sections
}

func (fb2 *FB2) BodySectionTitles() []string {
	var names []string
	for _, body := range fb2.Bodies {
		for _, section := range body.Sections {
			for _, paragraph := range section.Title.Paragraphs {
				names = append(names, paragraph.Text)
			}
		}
	}
	return names
}

var (
	ErrMissionOrIncorrectXMLNS      = errors.New("missing or incorrect xmlns")
	ErrMissionOrIncorrectXMLNSXLink = errors.New("missing or incorrect xmlns:xlink")
	ErrMissionOrIncorrectXMLNSGenre = errors.New("missing or incorrect xmlns:genre")
	ErrMissionOrIncorrectXMLNSXS    = errors.New("missing or incorrect xmlns:xs")
)

func (fb2 *FB2) ValidateNamespaces() error {
	var errs []error
	if fb2.XMLNS == nil || *fb2.XMLNS != XMLNS_FICTIONBOOK {
		errs = append(errs, ErrMissionOrIncorrectXMLNS)
	}
	if fb2.XMLNSXLink == nil || *fb2.XMLNSXLink != XMLNS_XLINK {
		errs = append(errs, ErrMissionOrIncorrectXMLNSXLink)
	}
	if fb2.XMLNSGenre == nil || *fb2.XMLNSGenre != XMLNS_GENRE {
		errs = append(errs, ErrMissionOrIncorrectXMLNSGenre)
	}
	if fb2.XMLNSXS == nil || *fb2.XMLNSXS != XMLNS_XS {
		errs = append(errs, ErrMissionOrIncorrectXMLNSXS)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (fb2 *FB2) FixNamespaces() {
	if err := fb2.ValidateNamespaces(); err != nil {
		if errors.Is(err, ErrMissionOrIncorrectXMLNS) {
			if fb2.XMLNS == nil {
				fb2.XMLNS = new(string)
			}
			*fb2.XMLNS = XMLNS_FICTIONBOOK
		}
		if errors.Is(err, ErrMissionOrIncorrectXMLNSXLink) {
			if fb2.XMLNSXLink == nil {
				fb2.XMLNSXLink = new(string)
			}
			*fb2.XMLNSXLink = XMLNS_XLINK
		}
		if errors.Is(err, ErrMissionOrIncorrectXMLNSGenre) {
			if fb2.XMLNSGenre == nil {
				fb2.XMLNSGenre = new(string)
			}
			*fb2.XMLNSGenre = XMLNS_GENRE
		}
		if errors.Is(err, ErrMissionOrIncorrectXMLNSXS) {
			if fb2.XMLNSXS == nil {
				fb2.XMLNSXS = new(string)
			}
			*fb2.XMLNSXS = XMLNS_XS
		}
	}
}
