package gofb2_test

import (
	"os"
	"testing"

	gofb2 "github.com/0x0FACED/go-fb2/v2"
)

func getTestFB2() *gofb2.FB2 {
	xmlns := gofb2.XMLNS_FICTIONBOOK
	xmlnsXLink := gofb2.XMLNS_XLINK
	xmlnsGenre := gofb2.XMLNS_GENRE
	xmlnsXS := gofb2.XMLNS_XS

	return &gofb2.FB2{
		XMLNS:      &xmlns,
		XMLNSXLink: &xmlnsXLink,
		XMLNSGenre: &xmlnsGenre,
		XMLNSXS:    &xmlnsXS,
		Stylesheet: []string{"stylesheet.css"},
		Description: gofb2.Description{
			TitleInfo: gofb2.TitleInfo{
				Genre: []string{"science-fiction", "fantasy"},
				Author: []gofb2.AuthorType{
					{
						FirstName:  "John",
						LastName:   "Doe",
						MiddleName: "A.",
					},
				},
				BookTitle: "Test Book",
				Annotation: gofb2.Annotation{
					Paragraphs: []gofb2.Paragraph{
						{Text: "This is a test annotation."},
					},
				},
				Coverpage: gofb2.Coverpage{
					Image: &gofb2.Image{Href: "#cover"},
				},
				Keywords: "test, example",
				Date: gofb2.Date{
					Value: "2024-12-01",
					Text:  "1 December 2024",
				},
				Lang:    "en",
				SrcLang: "ru",
			},
			DocumentInfo: gofb2.DocumentInfo{
				Author: []gofb2.AuthorType{
					{
						FirstName: "Editor",
						LastName:  "Smith",
					},
				},
				ProgramUsed: "Go FB2 Parser",
				Date: gofb2.Date{
					Value: "2024-12-01",
					Text:  "1 December 2024",
				},
				ID:      "1234567890",
				Version: 1.0,
				History: "Initial creation.",
			},
		},
		Bodies: []gofb2.Body{
			{
				Name: "Main",
				Sections: []gofb2.Section{
					{
						Title: gofb2.Title{
							Paragraphs: []gofb2.Paragraph{
								{Text: "Chapter 1"},
							},
						},
						Paragraphs: []gofb2.Paragraph{
							{Text: "This is the first paragraph of the test book."},
							{Text: "This is the second paragraph."},
						},
					},
				},
			},
		},
		Binary: []gofb2.Binary{
			{
				ID:          "cover",
				ContentType: "image/jpeg",
				Value:       "BASE64_ENCODED_IMAGE_DATA",
			},
		},
	}
}

func TestUnmarshal(t *testing.T) {
	parser := &gofb2.Parser{}

	tests := []struct {
		name      string
		filePath  string
		isStrict  bool
		expectErr bool
	}{
		{
			name:      "Test Unmarshal test_input.fb2",
			filePath:  "../test_input.fb2",
			isStrict:  true,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open(tt.filePath)
			if err != nil {
				t.Fatalf("Failed to open file %s: %v", tt.filePath, err)
			}
			defer file.Close()

			result, err := parser.Unmarshal(file, tt.isStrict)
			if (err != nil) != tt.expectErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectErr, err)
			}

			if result == nil {
				t.Errorf("Expected non-nil result for %s", tt.name)
			}
		})
	}
}

func TestUnmarshalFromFile(t *testing.T) {
	parser := &gofb2.Parser{}

	tests := []struct {
		name      string
		filePath  string
		isStrict  bool
		expectErr bool
	}{
		{
			name:      "Test UnmarshalFromFile test_input.fb2",
			filePath:  "../test_input.fb2",
			isStrict:  false,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.UnmarshalFromFile(tt.filePath, tt.isStrict)
			if (err != nil) != tt.expectErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectErr, err)
			}

			if result == nil {
				t.Errorf("Expected non-nil result for %s", tt.name)
			}
		})
	}
}

func TestMarshal(t *testing.T) {
	parser := &gofb2.Parser{}

	testFB2 := getTestFB2()
	tests := []struct {
		name      string
		fb2       *gofb2.FB2
		pretty    bool
		expectErr bool
	}{
		{
			name:      "Test Marshal with pretty formatting",
			fb2:       testFB2,
			pretty:    true,
			expectErr: false,
		},
		{
			name:      "Test Marshal without pretty formatting",
			fb2:       testFB2,
			pretty:    false,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.Marshal(tt.fb2, tt.pretty)
			if (err != nil) != tt.expectErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectErr, err)
			}

			if result == nil {
				t.Errorf("Expected non-nil result for %s", tt.name)
			}
		})
	}
}

func TestMarshalToFile(t *testing.T) {
	parser := &gofb2.Parser{}

	testFB2 := getTestFB2()
	tests := []struct {
		name      string
		fb2       *gofb2.FB2
		filePath  string
		pretty    bool
		expectErr bool
	}{
		{
			name:      "Test MarshalToFile with pretty formatting",
			fb2:       testFB2,
			filePath:  "TestMarshalToFile_output_pretty.fb2",
			pretty:    true,
			expectErr: false,
		},
		{
			name:      "Test MarshalToFile without pretty formatting",
			fb2:       testFB2,
			filePath:  "TestMarshalToFile_output_no_pretty.fb2",
			pretty:    false,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parser.MarshalToFile(tt.fb2, tt.filePath, tt.pretty)
			if (err != nil) != tt.expectErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectErr, err)
			}

			if _, err := os.Stat(tt.filePath); err != nil {
				t.Errorf("Failed to find output file %s: %v", tt.filePath, err)
			}
		})
	}
}

func TestLoadFrimFileAndSaveToFile(t *testing.T) {
	parser := &gofb2.Parser{}

	testFB2, err := parser.UnmarshalFromFile("../test_input.fb2", true)
	expectErr := false
	if (err != nil) != expectErr {
		t.Errorf("Expected error: %v, got: %v", expectErr, err)
	}

	err = parser.MarshalToFile(testFB2, "TestLoadFrimFileAndSaveToFile_output.fb2", true)
	if (err != nil) != expectErr {
		t.Errorf("Expected error: %v, got: %v", expectErr, err)
	}
}
