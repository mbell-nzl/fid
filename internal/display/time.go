package display

import "time"

// Clock shows a timestamp in now's time zone: HH:MM today, with the date
// on any other day. JSON retains the UTC timestamp.
func Clock(stamp string, now time.Time) string {
	at, err := time.Parse(time.RFC3339, stamp)
	if err != nil {
		return stamp
	}
	at = at.In(now.Location())
	if at.Format("2006-01-02") != now.Format("2006-01-02") {
		return at.Format("2006-01-02 15:04")
	}
	return at.Format("15:04")
}
