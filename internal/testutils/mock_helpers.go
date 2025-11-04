package testutils

// MockTracker provides mock tracking functionality
type MockTracker struct {
	CallCount   map[string]int
	CallHistory []string
	LastCommand string
}

// Track records a method call
func (mt *MockTracker) Track(methodName string) {
	mt.CallCount[methodName]++
	mt.CallHistory = append(mt.CallHistory, methodName)
	mt.LastCommand = methodName
}

// GetCallCount returns the number of times a method was called
func (mt *MockTracker) GetCallCount(methodName string) int {
	return mt.CallCount[methodName]
}

// Reset clears all tracking data
func (mt *MockTracker) Reset() {
	mt.CallCount = make(map[string]int)
	mt.CallHistory = make([]string, 0)
	mt.LastCommand = ""
}
