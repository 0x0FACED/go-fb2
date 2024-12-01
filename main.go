package main

import (
	"fmt"
	"log"

	gofb2 "github.com/0x0FACED/go-fb2/v2"
)

func main() {
	file := "/home/podliva/Documents/test-fb2.fb2"

	parser := gofb2.NewParser()

	// Парсинг FB2 файла
	fb2Book, err := parser.UnmarshalFromFile(file, true)
	if err != nil {
		log.Fatalf("Ошибка чтения FB2: %v", err)
	}

	fmt.Printf("Книга загружена: %+v\n", fb2Book)

	parser.MarshalToFile(fb2Book, "output.fb2", true)

		
}
