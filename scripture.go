package orthocal

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
)

// These books have a single chapter
var chapterlessBooks = map[string]bool{
	"Obad":  true,
	"Phlm":  true,
	"2John": true,
	"3John": true,
	"Jude":  true,
}

type Verse struct {
	Book    string `json:"book"`
	Chapter uint16 `json:"chapter"`
	Verse   uint16 `json:"verse"`
	Content string `json:"content"`
}

type Passage []Verse

type Bible struct {
	db *sql.DB
}

var (
	refRe   = regexp.MustCompile(`(?:([\w\s]+)\s+)?(\d.*)`)
	rangeRe = regexp.MustCompile(`(?:(\d+)[\.:])?(\d+)(?:-(?:(\d+)[\.:])?(\d+))?`)

	ReferenceParseError = errors.New("Error parsing scripture reference")
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
	Prov 10, 3, 8

	NOTE: this function directly interpolates values from the reference into
	SQL. This is safe as long as the provided reference is coming from the
	database. In other words, this method might be unsafe when used to lookup
	user-provided references.
*/
func (self *Bible) Lookup(reference string) Passage {
	var passage Passage

	sql, e := self.buildSQL(reference)
	if e != nil {
		log.Printf("Error building scripture lookup SQL '%s': %+v.", reference, e)
		return passage
	}

	rows, e := self.db.Query(sql)
	if e != nil {
		log.Printf("Got error querying the database for scripture '%s': %#n.", reference, e)
		return passage
	}
	defer rows.Close()

	for rows.Next() {
		var verse Verse
		rows.Scan(&verse.Book, &verse.Chapter, &verse.Verse, &verse.Content)
		passage = append(passage, verse)
	}

	return passage
}

func (self *Bible) buildSQL(reference string) (string, error) {
	var conditionals []string

	sql := "select book, chapter, verse, content from bible\n"

	book := ""
	for _, passage := range regexp.MustCompile(`\s*;\s*`).Split(reference, -1) {
		// Get the book and specification
		groups := refRe.FindStringSubmatch(passage)
		if len(groups) < 3 {
			return "", ReferenceParseError
		}

		// If the book is specified in this group, use it, else default to the previous one
		if len(groups[1]) > 0 {
			book = strings.Replace(groups[1], " ", "", -1)
		}
		specification := groups[2]

		// Create a conditional for each verse range in the specification
		var previousChapter string
		for _, verseRange := range regexp.MustCompile(`,\s*`).Split(specification, -1) {
			var conditional, chapter string

			m := rangeRe.FindStringSubmatch(verseRange)
			if len(m) < 5 {
				return "", ReferenceParseError
			}

			if _, ok := chapterlessBooks[book]; ok {
				// single-chapter books
				chapter = "1"
			} else {
				// multi-chapter books
				if len(m[1]) > 0 {
					chapter = m[1]
				} else {
					// The chapter is implicit, so we use the chapter from the
					// previous range if there is no chapter from the previous
					// range, then this range is specifying a full chapter and
					// the chapter number will be in m[2] and chapter will be
					// the empty string.
					chapter = previousChapter
				}
			}

			if len(m[4]) > 0 {
				if len(m[3]) > 0 && m[3] != m[1] {
					// Handle ranges that span chapters
					conditional = fmt.Sprintf(`(book = "%s" and ((chapter = %s and verse >= %s) or (chapter = %s and verse <= %s)))`, book, chapter, m[2], m[3], m[4])
				} else {
					// Handle ranges that are contained within a single chapter
					conditional = fmt.Sprintf(`(book = "%s" and chapter = %s and verse between %s and %s)`, book, chapter, m[2], m[4])
				}
			} else if len(chapter) > 0 {
				// Handle a single verse
				conditional = fmt.Sprintf(`(book = "%s" and chapter = "%s" and verse = "%s")`, book, chapter, m[2])
			} else {
				// Here we handle full chapters
				conditional = fmt.Sprintf(`(book = "%s" and chapter = "%s")`, book, m[2])
			}

			conditionals = append(conditionals, conditional)

			// Remember the most recently used chapter, unless it was a full chapter
			// If it was a full chapter, it can't be reused in a subsequent range.
			if len(m[3]) > 0 {
				previousChapter = m[3]
			} else if len(m[1]) > 0 {
				previousChapter = m[1]
			}
		}
	}

	if len(conditionals) > 1 {
		// "Or" all the conditional clauses together
		sql += "where (\n" + strings.Join(conditionals, "\nor ") + "\n)\n"
	} else {
		sql += "where " + conditionals[0] + "\n"
	}

	return sql, nil
}
