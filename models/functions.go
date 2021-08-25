package models

import (
	"fmt"
	"strconv"
	"time"
)

func (p *SearchResponse) AddDays(inputTime time.Time, days int64) time.Time {
	fmt.Println("adding days ", days)
	return inputTime.AddDate(0, 0, int(days))
}

func TransformDate(date time.Time) string {
	dateForKey := strconv.Itoa(date.Day()) + date.Month().String() + strconv.Itoa(date.Year())
	fmt.Println("date ", dateForKey)
	return dateForKey
}
