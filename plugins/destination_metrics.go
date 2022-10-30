package plugins

type DestinationMetrics struct {
	// Errors number of errors / failed writes
	Errors uint64
	// Writes number of successful writes
	Writes uint64
}
