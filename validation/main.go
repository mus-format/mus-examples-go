package main

import "github.com/ymz-ncnk/assert"

func init() {
	assert.On = true
}

// Shows how you can validate different data types during unmarshaling. Check
// out the functions to get more information.
func main() {
	ValidateString()
	ValidateSlice()
	ValidateMap()
	ValidateStruct()
}
