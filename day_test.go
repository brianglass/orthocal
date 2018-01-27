package orthocal_test

import (
	"orthocal"
	"testing"
)

func TestDB(t *testing.T) {
	orthocal.TestDB()
	t.Fail()
}
