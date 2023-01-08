package service

import (
	"time"
)

func GetFirstAndLastDayFromMonth(month string) (time.Time, time.Time, error) {
	date, err := time.Parse("2006-01", month)
	if err != nil {
		return time.Now(), time.Now(), err
	}
	currentYear, currentMonth, _ := date.Date()
	currentLocation := date.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	return firstOfMonth, lastOfMonth, nil

}
