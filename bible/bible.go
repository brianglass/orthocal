package main

import (
	"database/sql"
	"encoding/xml"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Passage struct {
	book    string
	chapter int
	verse   int
	content string
}

func main() {
	xmlFile, _ := os.Open("eng-kjv_usfx.xml")
	defer xmlFile.Close()

	db, e := sql.Open("sqlite3", "kjv.db")
	if e != nil {
		log.Printf("Got error opening database: %+v.", e)
		os.Exit(1)
	}
	defer db.Close()

	_, e = db.Exec(`CREATE TABLE bible (book, chapter INTEGER, verse INTEGER, content TEXT, PRIMARY KEY(book, chapter, verse))`)
	if e != nil {
		log.Printf("Got error creating table: %+v.", e)
		os.Exit(1)
	}

	c := make(chan Passage)
	go Parse(xmlFile, c)
	for passage := range c {
		_, e = db.Exec(`INSERT INTO bible VALUES (?, ?, ?, ?)`, passage.book, passage.chapter, passage.verse, passage.content)
		if e != nil {
			log.Printf("Got error inserting row: %+v.", e)
			os.Exit(1)
		}
	}

	// Inserts are slightly faster if we do this after the table is built
	_, e = db.Exec(`CREATE INDEX bible_book_chapter ON bible(book, chapter)`)
	if e != nil {
		log.Printf("Got error creating index: %+v.", e)
		os.Exit(1)
	}
}

func Parse(reader io.Reader, c chan Passage) {
	var book, content string
	var chapter, verse int
	var isValidContent bool

	decoder := xml.NewDecoder(reader)

	for {
		token, e := decoder.Token()

		if e != nil {
			if e != io.EOF {
				log.Printf("Error parsing usfx file: %+v", e)
			}
			close(c)
			return
		}

		switch element := token.(type) {
		case xml.StartElement:
			switch element.Name.Local {
			case "book":
				book = getAttribute(element, "id")
			case "c":
				chapter, _ = strconv.Atoi(getAttribute(element, "id"))
			case "v":
				verse, _ = strconv.Atoi(getAttribute(element, "id"))
				isValidContent = true
			case "ve":
				// Verse end
				c <- Passage{book, chapter, verse, strings.TrimSpace(content)}
				content = ""
				isValidContent = false
			case "f":
				// ignore footnotes
				isValidContent = false
			}
		case xml.EndElement:
			switch element.Name.Local {
			case "f":
				// Exclude this text from the verse content
				isValidContent = true
			}
		case xml.CharData:
			if isValidContent {
				content = strings.TrimSpace(content)
				text := strings.TrimSpace(string(element))
				text = strings.Replace(text, "¶", "", -1)
				first, _ := utf8.DecodeRuneInString(text)
				if first == '.' || first == ',' || first == '?' || first == '!' || first == ':' || first == ';' {
					content += text
				} else {
					content += " " + text
				}
			}
		case xml.ProcInst, xml.Directive, xml.Comment:
		}

	}
}

func getAttribute(element xml.StartElement, attribute string) string {
	for i := range element.Attr {
		if element.Attr[i].Name.Local == attribute {
			return element.Attr[i].Value
		}
	}

	return ""
}
