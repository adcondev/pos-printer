package test

// FakeCapabilityBase provides sharedcommands fake functionality
type FakeCapabilityBase struct {
	Buffer   []byte
	Tracker  *MockTracker
	StateMap map[string]interface{}
}

// NewFakeCapabilityBase creates a new fake base
func NewFakeCapabilityBase() *FakeCapabilityBase {
	return &FakeCapabilityBase{
		Buffer:   make([]byte, 0),
		Tracker:  NewMockTracker(),
		StateMap: make(map[string]interface{}),
	}
}

// AppendCommand adds a command to the buffer and tracks it
func (f *FakeCapabilityBase) AppendCommand(cmd []byte, methodName string) {
	f.Buffer = append(f.Buffer, cmd...)
	f.Tracker.Track(methodName)
}

// GetBuffer returns the accumulated buffer
func (f *FakeCapabilityBase) GetBuffer() []byte {
	return f.Buffer
}

// GetState returns the current state map
func (f *FakeCapabilityBase) GetState() map[string]interface{} {
	return f.StateMap
}

// SetState updates a state value
func (f *FakeCapabilityBase) SetState(key string, value interface{}) {
	f.StateMap[key] = value
}

// Reset clears all data
func (f *FakeCapabilityBase) Reset() {
	f.Buffer = make([]byte, 0)
	f.Tracker.Reset()
	f.StateMap = make(map[string]interface{})
}
