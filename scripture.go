package orthocal

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"
)

type Verse struct {
	Chapter uint16 `json:"chapter"`
	Verse   uint16 `json:"verse"`
	Content string `json:"content"`
}

type Passage []Verse

type Bible struct {
	db *sql.DB
}

var (
	refRe   = regexp.MustCompile(`([\w\s]+)\s+(\d.*)`)
	rangeRe = regexp.MustCompile(`(?:(\d+)\.)?(\d+)(?:-(?:(\d+)\.)?(\d+))?`)
)

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
	Mark 15.22, 25, 33-41
	1 John 2.7-17
	Jude 1-10
	1 Cor 5.6-8; Gal 3.13-14

	NOTE: this function directly interpolates values from the reference into
	SQL. This is safe as long as the provided reference is coming from the
	database. In other words, this method might be unsafe when used to lookup
	user-provided references.
*/
func (self *Bible) Lookup(reference string) Passage {
	var passage Passage

	sql := self.buildSQL(reference)
	rows, e := self.db.Query(sql)
	if e != nil {
		log.Printf("Got error querying the database for scripture '%s': %#n.", reference, e)
		return passage
	}
	defer rows.Close()

	for rows.Next() {
		var verse Verse
		rows.Scan(&verse.Chapter, &verse.Verse, &verse.Content)
		passage = append(passage, verse)
	}

	return passage
}

func (self *Bible) buildSQL(reference string) string {
	var conditionals []string

	sql := "select chapter, verse, content from bible\n"

	for _, passage := range regexp.MustCompile(`\s*;\s*`).Split(reference, 4) {
		// Get the book and specification
		groups := refRe.FindStringSubmatch(passage)
		book := strings.Replace(groups[1], " ", "", -1)
		specification := groups[2]

		// Create a conditional for each verse range in the specification
		chapter := "1"
		for _, verseRange := range regexp.MustCompile(`,\s*`).Split(specification, 4) {
			var conditional string

			m := rangeRe.FindStringSubmatch(verseRange)

			defaultChapter := m[1]
			if len(defaultChapter) == 0 {
				// The chapter is implicit, so we use the chapter from the previous range
				defaultChapter = chapter
			}

			if len(m[4]) > 0 {
				if len(m[3]) > 0 && m[3] != m[1] {
					// Handle ranges that span chapters
					conditional = fmt.Sprintf(`(book = "%s" and ((chapter = %s and verse >= %s) or (chapter = %s and verse <= %s)))`, book, defaultChapter, m[2], m[3], m[4])
				} else {
					// Handle ranges that are contained within a single chapter
					conditional = fmt.Sprintf(`(book = "%s" and chapter = %s and verse between %s and %s)`, book, defaultChapter, m[2], m[4])
				}
			} else {
				// Handle a single verse
				conditional = fmt.Sprintf(`(book = "%s" and chapter = "%s" and verse = "%s")`, book, defaultChapter, m[2])
			}

			conditionals = append(conditionals, conditional)

			// Remember the most recently used chapter
			if len(m[3]) > 0 || len(m[1]) > 0 {
				chapter = m[3]
				if len(chapter) == 0 {
					chapter = m[1]
				}
			}
		}
	}

	if len(conditionals) > 1 {
		// "Or" all the conditional clauses together
		sql += "where (\n" + strings.Join(conditionals, "\nor ") + "\n)\n"
	} else {
		sql += "where " + conditionals[0] + "\n"
	}

	return sql
}
