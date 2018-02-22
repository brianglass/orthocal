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

/*
	Parse a scripture reference and return the scripture passage.  A reference
	includes a book name and a verse specification.  A specification can
	include multiple verse ranges separated by commas, where a range is a
	collection of 1 or more contiguous verses. Chapter specification can be
	implicit. A reference can look like any of the following:

	Matt 1.1-25
	Matt 4.25-5.13
	Matt 10.32-36, 11.1
	Matt 6.31-34, 7.9-11
	Matt 10.1, 5-8

	NOTE: this function directly interpolates values from the reference into
	SQL. This is safe as long as the provided reference is coming from the
	database. In other words, this method might be unsafe when used to lookup
	user-provided references.
*/
func (self *Bible) Lookup(reference string) Passage {
	var passage Passage

	sql := self.convertReferenceToSQL(reference)
	rows, e := self.db.Query(sql)
	defer rows.Close()

	if e != nil {
		log.Printf("Got error querying the database: %#n.", e)
		return passage
	}

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

	sql := "select book, chapter, verse, content\nfrom bible\n"

	// Break out book and range specifications
	groups := regexp.MustCompile(`\s+`).Split(reference, 2)
	book, specification := groups[0], groups[1]

	sql += fmt.Sprintf("where book = \"%s\"\n", book)

	// Create a conditional for each verse range in the specification
	for _, verseRange := range regexp.MustCompile(`,\s*`).Split(specification, 4) {
		var conditional string

		m := refRe.FindStringSubmatch(verseRange)

		defaultChapter := m[1]
		if len(defaultChapter) == 0 {
			// The chapter is implicit, so we use the chapter from the previous range
			defaultChapter = chapter
		}

		if len(m[4]) > 0 {
			if len(m[3]) > 0 && m[3] != m[1] {
				// Handle ranges that span chapters
				conditional = fmt.Sprintf(`((chapter = %s and verse >= %s) or (chapter = %s and verse <= %s))`, defaultChapter, m[2], m[3], m[4])
			} else {
				// Handle ranges that are contained within a single chapter
				conditional = fmt.Sprintf(`(chapter = %s and verse between %s and %s)`, defaultChapter, m[2], m[4])
			}
		} else {
			// Handle a single verse
			conditional = fmt.Sprintf(`(chapter = "%s" and verse = "%s")`, defaultChapter, m[2])
		}

		conditionals = append(conditionals, conditional)

		// Remember the most recently used chapter
		chapter = m[3]
		if len(chapter) == 0 {
			chapter = m[1]
		}
	}

	if len(conditionals) > 1 {
		// "Or" all the conditional clauses together
		sql += "and (\n" + strings.Join(conditionals, "\nor ") + "\n)\n"
	} else {
		sql += "and " + conditionals[0] + "\n"
	}

	return sql
}
