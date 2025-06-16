package main

// IsTestMode returns true if the code is running in test mode
var inTestMode = false

// SetTestMode enables test mode for testing
func SetTestMode(enabled bool) {
	inTestMode = enabled
}

// IsTestMode returns true if the code is running in test mode
func IsTestMode() bool {
	return inTestMode
}
