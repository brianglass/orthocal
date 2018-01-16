package orthocal

import (
	"errors"
	"time"
)

// Pascha functions

func ComputeJulianPascha(year int) (int, int) {
	// Use the Meeus Julian algorithm to calculate the Julian date
	// See https://en.wikipedia.org/wiki/Computus#Meeus'_Julian_algorithm
	a := year % 4
	b := year % 7
	c := year % 19
	d := (19*c + 15) % 30
	e := (2*a + 4*b - d + 34) % 7
	month := (d + e + 114) / 31
	day := (d+e+114)%31 + 1
	return month, day
}

func ComputeJulianDayPascha(year int) int {
	month, day := ComputeJulianPascha(year)
	return JulianDateToJulianDay(year, month, day)
}

func ComputeGregorianPascha(year int) (time.Time, error) {
	month, day := ComputeJulianPascha(year)
	gregorianDate, e := JulianToGregorian(year, month, day)
	if e != nil {
		return time.Now(), e
	}
	return gregorianDate, nil
}

// Conversion functions

func JulianToGregorian(year, month, day int) (time.Time, error) {
	// This will be incorrect outside the range 2001-2099 for 2 reasons:
	// 1. The offset of 13 is incorrect outside the range 1900-2099.
	// 2. if the Julian date is in February and on a year that is divisible by
	//    100, the Go time module will incorrectly add the offset because these years
	//    are leap years on the Julian, but not on the Gregorian.
	if year < 2001 || year > 2099 {
		return time.Now(), errors.New("The year must be between 1900 and 2099")
	}

	// Add an offset of 13 to convert from Julian to Gregorian
	julianDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	return julianDate.AddDate(0, 0, 13), nil
}

func JulianDateToJulianDay(year, month, day int) int {
	// See https://en.wikipedia.org/wiki/Julian_day#Converting_Julian_calendar_date_to_Julian_Day_Number
	return 367*year - (7*(year+5001+(month-9)/7))/4 + (275*month)/9 + day + 1729777
}

func GregorianDateToJulianDay(gregorianDate time.Time) int {
	// This function mimic's PHP's gregoriantojd()

	// month is an integer from 1-12
	// day is an integer from 1-31
	// year is an integer from -4714 and 9999

	year, month, day := gregorianDate.Date()

	if month > 2 {
		month -= 3
	} else {
		month += 9
		year--
	}

	// break up the year into the leftmost 2 digits (century) and the rightmost 2 digits
	century := year / 100
	ya := year - 100*century

	return (146097*century)/4 + (1461*ya)/4 + (153*int(month)+2)/5 + day + 1721119
}
