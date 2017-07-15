package ght

import "log"

// OptionalLogger only logs when the verbose variable is true.
type OptionalLogger struct {
	verbose bool
}

// New creates a new VerboseLogger using the referenced bool to set verbosity.
func (l *OptionalLogger) New(v *bool) {
	l.verbose = *v
}

// IsVerbose specifies if the logger should log human friendly output. This is nice for a program that has either human friendly verbose output or terse machine friendly return codes.
func (l *OptionalLogger) IsVerbose() bool {
	return l.verbose
}

func (l *OptionalLogger) Write(p []byte) (n int, err error) {
	if l.verbose {
		log.Print(string(p[:len(p)]))
	}
	return 0, nil
}
