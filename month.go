package recurrence

import (
	"encoding/json"
	"fmt"
	"time"
)

// A Month represents a month of the year. Just like time.Month.
type Month time.Month

const (
	January Month = 1 + iota
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
)

func (self Month) IsOccurring(t time.Time) bool {
	return t.Month() == time.Month(self)
}

func (self Month) Occurrences(t TimeRange) chan time.Time {
	return t.occurrencesOfSchedule(self)
}

func (self Month) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{"month": time.Month(self).String()})
}

func (self *Month) UnmarshalJSON(b []byte) error {
	var err error

	switch string(b) {
	case `1`, `"January"`:
		*self = January
	case `2`, `"February"`:
		*self = February
	case `3`, `"March"`:
		*self = March
	case `4`, `"April"`:
		*self = April
	case `5`, `"May"`:
		*self = May
	case `6`, `"June"`:
		*self = June
	case `7`, `"July"`:
		*self = July
	case `8`, `"August"`:
		*self = August
	case `9`, `"September"`:
		*self = September
	case `10`, `"October"`:
		*self = October
	case `11`, `"November"`:
		*self = November
	case `12`, `"December"`:
		*self = December
	default:
		*self = 0
		err = fmt.Errorf("Weekday cannot unmarshal %s", b)
	}
	return err
}
