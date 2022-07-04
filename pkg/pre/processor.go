package pre

import "io"

// Processor is a pre-processor interface
type Processor interface {
	Process()
	// GetSource() io.Reader
	SetSource(io.Reader)
	// GetTarget() io.Writer
	SetTarget(io.Writer)
}
