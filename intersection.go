package recurrence

import (
	"encoding/json"
	"fmt"
	"time"
)

// Computes the set intersection of a slice of Schedules.
type Intersection []Schedule

// Implement Schedule interface.
func (self Intersection) IsOccurring(t time.Time) bool {
	for _, r := range self {
		if !r.IsOccurring(t) {
			return false
		}
	}

	return true
}

// Implement Schedule interface.
func (self Intersection) Occurrences(t TimeRange) []time.Time {
	done := make(chan bool, len(self))
	candidates := make(chan time.Time, 100)
	ts := make([]time.Time, 0)

	for _, schedule := range self {
		go func(schedule Schedule) {
			for _, oc := range schedule.Occurrences(t) {
				candidates <- oc
			}
			done <- true
		}(schedule)
	}

	candidatesMap := make(map[string]int)
	parallelDone := 0
	for parallelDone < len(self) {
		select {
		case selected := <-candidates:
			key := selected.Format("20060102")
			foundCount, _ := candidatesMap[key]
			newFoundCount := foundCount + 1
			candidatesMap[key] = newFoundCount
			if newFoundCount == len(self) {
				ts = append(ts, selected)
			}
		case <-done:
			parallelDone++
		}
	}

	// We can safely close channel done now
	close(done)

	// What if we somehow have some residual data in candidates channel?
	stillLopp := true
	for stillLopp {
		select {
		case selected := <-candidates:
			key := selected.Format("20060102")
			foundCount, _ := candidatesMap[key]
			newFoundCount := foundCount + 1
			candidatesMap[key] = newFoundCount
			if newFoundCount == len(self) {
				ts = append(ts, selected)
			}
			// This means there are no more elements left in our channel
		default:
			// We definitely don't have anything in this channel
			stillLopp = false
		}
	}

	// We can also safely close the candidates channel
	close(candidates)

	return ts
}

// Implement json.Marshaler interface.
func (self Intersection) MarshalJSON() ([]byte, error) {
	type faux Intersection
	return json.Marshal(struct {
		faux `json:"intersection"`
	}{faux: faux(self)})
}

// Implement json.Unmarshaler interface.
func (self *Intersection) UnmarshalJSON(b []byte) error {
	var mixed interface{}

	json.Unmarshal(b, &mixed)

	switch mixed.(type) {
	case []interface{}:
		for _, value := range mixed.([]interface{}) {
			bytes, _ := json.Marshal(value)
			schedule, err := ScheduleUnmarshalJSON(bytes)
			if err != nil {
				return err
			}
			*self = append(*self, schedule)
		}
	default:
		return fmt.Errorf("intersection must be a slice")
	}

	return nil
}
