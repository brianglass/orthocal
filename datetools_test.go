package orthocal_test

import (
	"github.com/brianglass/orthocal"
	"testing"
	"time"
)

var fixture_gregorian_pascha = []time.Time{
	time.Date(2008, 4, 27, 0, 0, 0, 0, time.Local),
	time.Date(2009, 4, 19, 0, 0, 0, 0, time.Local),
	time.Date(2010, 4, 4, 0, 0, 0, 0, time.Local),
	time.Date(2011, 4, 24, 0, 0, 0, 0, time.Local),
}

var fixture_julian_pascha = []time.Time{
	time.Date(2008, 4, 14, 0, 0, 0, 0, time.Local),
	time.Date(2009, 4, 6, 0, 0, 0, 0, time.Local),
	time.Date(2010, 3, 22, 0, 0, 0, 0, time.Local),
	time.Date(2011, 4, 11, 0, 0, 0, 0, time.Local),
}

func TestGregorianDateToJDN(t *testing.T) {
	actual := orthocal.GregorianDateToJDN(2018, 1, 15)
	expected := 2458134
	if actual != expected {
		t.Fatalf("GregorianDateToJDN should have returned %d but returned %d", expected, actual)
	}

	actual = orthocal.GregorianDateToJDN(2000, 5, 29)
	expected = 2451694
	if actual != expected {
		t.Fatalf("GregorianDateToJDN should have returned %d but returned %d", expected, actual)
	}
}

func TestComputeGregorianPascha(t *testing.T) {
	for _, expectedTime := range fixture_gregorian_pascha {
		pascha, e := orthocal.ComputeGregorianPascha(expectedTime.Year())
		if e != nil {
			t.Errorf("CalculateGregorianPascha had an error: %#v", e)
		}
		if pascha != expectedTime {
			t.Errorf("CalculateGregorianPascha should have returned %s but returned %s", expectedTime, pascha)
		}
	}
}

func TestComputeGregorianPaschaInvalid(t *testing.T) {
	expectedTime := time.Date(2100, 5, 2, 0, 0, 0, 0, time.Local)
	_, e := orthocal.ComputeGregorianPascha(expectedTime.Year())
	if e == nil {
		t.Errorf("CalculateGregorianPacha should return an error when dates are out of range")
	}
}

func TestComputeJulianPascha(t *testing.T) {
	for _, expectedTime := range fixture_julian_pascha {
		month, day := orthocal.ComputeJulianPascha(expectedTime.Year())
		if time.Month(month) != expectedTime.Month() {
			t.Errorf("CalculateJulianPascha should have returned %d for month but returned %d", month, expectedTime.Month())
		}
		if day != expectedTime.Day() {
			t.Errorf("CalculateJulianPascha should have returned %d for day but returned %d", day, expectedTime.Day())
		}
	}
}

func TestConvertJulianToGregorian(t *testing.T) {
	expected := time.Date(2008, 4, 27, 0, 0, 0, 0, time.Local)
	actual, e := orthocal.JulianToGregorian(2008, 4, 14)
	if e != nil {
		t.Fatalf("ConvertJulianToGregory return error: %#v", e)
	}

	if expected != actual {
		t.Errorf("ConvertJulianToGregory should have returned %v but returned %v", expected, actual)
	}

	expected = time.Date(2011, 4, 24, 0, 0, 0, 0, time.Local)
	actual, e = orthocal.JulianToGregorian(2011, 4, 11)
	if e != nil {
		t.Fatalf("ConvertJulianToGregory return error: %#v", e)
	}

	if expected != actual {
		t.Errorf("ConvertJulianToGregory should have returned %v but returned %v", expected, actual)
	}
}

func TestConvertJulianToGregorianInvalid(t *testing.T) {
	_, e := orthocal.JulianToGregorian(2100, 4, 14)
	if e == nil {
		t.Errorf("ConvertJulianToGregory should return an error when dates are out of range")
	}

	_, e = orthocal.JulianToGregorian(1900, 4, 14)
	if e == nil {
		t.Errorf("ConvertJulianToGregory should return an error when dates are out of range")
	}
}

func TestJulianDatetoJDN(t *testing.T) {
	expected := 2455676
	actual := orthocal.JulianDateToJDN(2011, 4, 11)
	if actual != expected {
		t.Errorf("JulianDateToJDN returned %d but should have returned %d", actual, expected)
	}
}

func TestComputJDNPascha(t *testing.T) {
	expected := 2455676
	actual := orthocal.ComputePaschaJDN(2011)
	if actual != expected {
		t.Errorf("ComputJDNPascha returned %d but should have returned %d", actual, expected)
	}
}

func TestComputePaschaDistance(t *testing.T) {
	distance, year := orthocal.ComputePaschaDistance(2018, 5, 9)
	expectedDistance, expectedYear := 31, 2018
	if distance != expectedDistance {
		t.Errorf("ComputePaschaDistance returned %d for the distance but should have returned %d", distance, expectedDistance)
	}
	if year != expectedYear {
		t.Errorf("ComputePaschaDistance returned %d for the year but should have returned %d", year, expectedYear)
	}

	distance, year = orthocal.ComputePaschaDistance(2018, 1, 1)
	expectedDistance, expectedYear = 260, 2017
	if distance != expectedDistance {
		t.Errorf("ComputePaschaDistance returned %d for the distance but should have returned %d", distance, expectedDistance)
	}
	if year != expectedYear {
		t.Errorf("ComputePaschaDistance returned %d for the year but should have returned %d", year, expectedYear)
	}
}

func TestWeekDayFromPDist(t *testing.T) {
	distance := 31
	expected := orthocal.Wednesday
	actual := orthocal.WeekDayFromPDist(distance)
	if actual != expected {
		t.Errorf("WeekDayFromPDist returned %d for the day but should have returned %d", actual, expected)
	}
}
