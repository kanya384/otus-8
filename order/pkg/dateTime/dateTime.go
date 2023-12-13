package dateTime

import (
	"fmt"
	"strconv"
	"strings"
)

type DateTime struct {
	hour   int
	minute int
}

var (
	ErrNotValidHourInterval   = fmt.Errorf("hour must be between 0 and 23")
	ErrNotValidMinuteInterval = fmt.Errorf("hour must be between 0 and 59")
)

func NewDateTime(hour int, minute int) (DateTime, error) {
	if hour < 0 && hour > 24 {
		return DateTime{}, ErrNotValidHourInterval
	}

	if minute < 0 && minute > 59 {
		return DateTime{}, ErrNotValidHourInterval
	}
	return DateTime{
		hour:   hour,
		minute: minute,
	}, nil
}
func NewDateTimeFromString(input string) (DateTime, error) {
	items := strings.Split(input, ":")
	if len(items) < 2 {
		return DateTime{}, fmt.Errorf("error parsing time")
	}
	hour, err := strconv.Atoi(items[0])
	if err != nil {
		return DateTime{}, err
	}

	minute, err := strconv.Atoi(items[1])
	if err != nil {
		return DateTime{}, err
	}
	return NewDateTime(hour, minute)
}

func (d DateTime) Hour() int {
	return d.hour
}

func (d DateTime) Minute() int {
	return d.minute
}

func (d DateTime) After(date DateTime) bool {
	if date.hour > d.hour {
		return false
	}

	if date.hour == d.hour && date.minute > d.hour {
		return false
	}

	return true
}

func (d DateTime) IsNil() bool {

	if d.hour == 0 && d.minute == 0 {
		return true
	}

	return false
}

func (d DateTime) ToString() string {
	var result string
	if d.hour < 10 {
		result = "0"
	}
	result += strconv.Itoa(d.hour) + ":"
	if d.minute < 10 {
		result += "0"
	}
	result += strconv.Itoa(d.minute)
	return result
}
