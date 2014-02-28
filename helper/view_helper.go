// These view helper methods provided extended functionality for building views
// using the Go template engine
package helper

import (
	"fmt"
	"html"
	"time"
)

type (
	// ViewHelper is passed into the view engine to access the
	// functionality extended by the object
	ViewHelper struct {
	}
)

// TrimString returns the specified number of characters from the
// specified string
func (this *ViewHelper) TrimString(s string, length int) string {
	if len(s) <= length {
		return html.UnescapeString(s)
	}

	return html.UnescapeString(s[0:length])
}

// UnescapeString can be used to display HTML in the view correctly
// and not as code
func (this *ViewHelper) UnescapeString(s string) string {
	return html.UnescapeString(s)
}

// ToInt converts a 64 bit float value to a displayble integer value
func (this *ViewHelper) ToInt(value float64) int {
	return int(value)
}

// ToMoney converts a 64 bit float value to a displayble money value
func (this *ViewHelper) ToMoney(value float64) string {
	return fmt.Sprintf("$%.2f", value)
}

// FormatDate converts a time.Time value to a formatted displayable datetime
func (this *ViewHelper) FormatDate(value time.Time) string {
	viewTime := value.Local().Format("2006-01-02 15:04:05")
	if viewTime == "0000-12-31 20:00:00" {
		return ""
	}

	return viewTime
}

// DurationToSeconds converts an integer representing time in milliseconds
// to fractions of a second.
func (this *ViewHelper) DurationToSeconds(timeMilli int) string {
	duration := float64(timeMilli) * .001 / 60
	return fmt.Sprintf("%.2f", duration)
}
