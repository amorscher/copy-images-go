// contains some utility functions for dates
package utils

import "time"

//RemoveMonths removes the numberOfMonthsToRemove from the provided date and returns a new one
func RemoveMonths(time time.Time, numberOfMonthsToRemove int) time.Time {
	return time.AddDate(0, numberOfMonthsToRemove*-1, 0)
}
