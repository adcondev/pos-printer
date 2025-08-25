package escpos

// MockPrinterCapability tracks method calls for testing and mocking printer commands
type MockPrinterCapability struct {
	// Track method calls
	TextCalled bool
	TextInput  string
	TextReturn []byte
	TextError  error

	PrintFeedAndPaperCalled bool
	PrintAndFeedPaperInput  byte
	PrintAndFeedPaperReturn []byte

	FormFeedCalled bool
	FormFeedReturn []byte

	PrintAndCarriageReturnCalled bool
	PrintAndCarriageReturnReturn []byte

	PrintAndLineFeedCalled bool
	PrintAndLineFeedReturn []byte
}

// Text mock
func (m *MockPrinterCapability) Text(n string) ([]byte, error) {
	m.TextCalled = true
	m.TextInput = n
	if m.TextError != nil {
		return nil, m.TextError
	}
	if m.TextReturn != nil {
		return m.TextReturn, nil
	}
	return []byte(n), nil
}

// PrintAndFeedPaper mock
func (m *MockPrinterCapability) PrintAndFeedPaper(n byte) []byte {
	m.PrintFeedAndPaperCalled = true
	m.PrintAndFeedPaperInput = n
	if m.PrintAndFeedPaperReturn != nil {
		return m.PrintAndFeedPaperReturn
	}
	return []byte{ESC, 'J', n}
}

// FormFeed mock
func (m *MockPrinterCapability) FormFeed() []byte {
	m.FormFeedCalled = true
	if m.FormFeedReturn != nil {
		return m.FormFeedReturn
	}
	return []byte{FF}
}

// PrintAndCarriageReturn mock
func (m *MockPrinterCapability) PrintAndCarriageReturn() []byte {
	m.PrintAndCarriageReturnCalled = true
	if m.PrintAndCarriageReturnReturn != nil {
		return m.PrintAndCarriageReturnReturn
	}
	return []byte{CR}
}

// PrintAndLineFeed mock
func (m *MockPrinterCapability) PrintAndLineFeed() []byte {
	m.PrintAndLineFeedCalled = true
	if m.PrintAndLineFeedReturn != nil {
		return m.PrintAndLineFeedReturn
	}
	return []byte{LF}
}

// MockLineSpacingCapability for testing line spacing commands
type MockLineSpacingCapability struct {
	SetLineSpacingCalled bool
	SetLineSpacingInput  byte
	SetLineSpacingReturn []byte

	SelectDefaultCalled bool
	SelectDefaultReturn []byte
}

func (m *MockLineSpacingCapability) SetLineSpacing(n byte) []byte {
	m.SetLineSpacingCalled = true
	m.SetLineSpacingInput = n
	if m.SetLineSpacingReturn != nil {
		return m.SetLineSpacingReturn
	}
	return []byte{ESC, '3', n}
}

func (m *MockLineSpacingCapability) SelectDefaultLineSpacing() []byte {
	m.SelectDefaultCalled = true
	if m.SelectDefaultReturn != nil {
		return m.SelectDefaultReturn
	}
	return []byte{ESC, '2'}
}
