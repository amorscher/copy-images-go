package utils_test

import (
	"copy-images/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRemoveMonthsRemovesTheNumberOfMonths(t *testing.T) {

	//GIVEN
	var givenDate time.Time
	givenDate, _ = time.Parse("2006-01-02", "2021-03-03")

	//WHEN
	var resultRemove2 = utils.RemoveMonths(givenDate, 2)
	var resultRemove1 = utils.RemoveMonths(givenDate, 1)

	//THEN
	assert.Equal(t, time.February, resultRemove1.Month(), "1 month in the past must be februrary")
	assert.Equal(t, time.January, resultRemove2.Month(), "2 month in the past must be january")

}

func TestRemoveMonthsRemovesTheNumberOfMonths_ChangingYear(t *testing.T) {

	//GIVEN
	var givenDate time.Time
	givenDate, _ = time.Parse("2006-01-02", "2021-03-03")

	//WHEN
	var resultRemove3 = utils.RemoveMonths(givenDate, 3)
	var resultRemove6 = utils.RemoveMonths(givenDate, 6)

	//THEN
	assert.Equal(t, time.December, resultRemove3.Month(), "3 month in the past must be January")
	assert.Equal(t, 2020, resultRemove3.Year(), "Year must be one year in the past")
	assert.Equal(t, time.September, resultRemove6.Month(), "6 month in the past must be october")
	assert.Equal(t, 2020, resultRemove6.Year(), "Year must be one year in the past")

}
