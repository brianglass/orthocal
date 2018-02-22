package orthocal

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"
)

type Verse struct {
	Book    string
	Chapter uint16
	Verse   uint16
	Content string
}

type Passage []Verse

type Bible struct {
	db *sql.DB
}

var refRe = regexp.MustCompile(`(?:(\d+)\.)?(\d+)(?:-(?:(\d+)\.)?(\d+))?`)

func NewBible(db *sql.DB) *Bible {
	var self Bible
	self.db = db
	return &self
}

func (self *Bible) Lookup(reference string) Passage {
	sql := self.convertReferenceToSQL(reference)
	rows, e := self.db.Query(sql)
	defer rows.Close()

	if e != nil {
		log.Printf("Got error querying the database: %#n.", e)
		return nil
	}

	var passage Passage
	for rows.Next() {
		var verse Verse
		rows.Scan(&verse.Book, &verse.Chapter, &verse.Verse, &verse.Content)
		passage = append(passage, verse)
	}

	return passage
}

func (self *Bible) convertReferenceToSQL(reference string) string {
	var conditionals []string
	var chapter string

	sql := "select book, chapter, verse, content from bible\n"
	groups := regexp.MustCompile(`\s+`).Split(reference, 2)
	book, specification := groups[0], groups[1]

	sql += fmt.Sprintf("where book = \"%s\" and (\n", book)

	references := regexp.MustCompile(`,\s*`).Split(specification, 4)
	for _, reference := range references {
		var conditional string
		m := refRe.FindStringSubmatch(reference)

		defaultChapter := m[1]
		if len(defaultChapter) == 0 {
			defaultChapter = chapter
		}

		if len(m[4]) > 0 {
			// Handle ranges
			if len(m[3]) > 0 && m[3] != m[1] {
				conditional = fmt.Sprintf(`((chapter = %s and verse >= %s) or (chapter = %s and verse <= %s))`, defaultChapter, m[2], m[3], m[4])
			} else {
				conditional = fmt.Sprintf(`(chapter = %s and verse between %s and %s)`, defaultChapter, m[2], m[4])
			}
		} else {
			// Handle a single verse
			conditional = fmt.Sprintf(`(chapter = "%s" and verse = "%s")`, defaultChapter, m[2])
		}

		conditionals = append(conditionals, conditional)

		chapter = m[3]
		if len(chapter) == 0 {
			chapter = m[1]
		}
	}

	sql += strings.Join(conditionals, "\nor ") + "\n)\n"

	return sql
}
