package utils_test

import (
	"copy-images/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThatItemCanBeFoundInArray(t *testing.T) {

	//GIVEN
	var testArray []string = []string{".png", ".jpeg", ".gif"}
	//WHEN
	var result = utils.ItemExists(testArray, ".jpeg")
	//THEN
	assert.True(t, result)

}

func TestThatItemCannotBeFoundInArray(t *testing.T) {

	//GIVEN
	var testArray []string = []string{".png", ".jpeg", ".gif"}
	//WHEN
	var result = utils.ItemExists(testArray, ".notFound")
	//THEN
	assert.False(t, result)

}
