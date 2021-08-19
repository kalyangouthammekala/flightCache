package models

import (
	"fmt"
	"time"
)

func (p *SearchResponse) AddDays(inputTime time.Time, days int64) time.Time {
	fmt.Println("adding days ", days)
	return inputTime.AddDate(0, 0, int(days))
}
