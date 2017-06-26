package ght

import "log"

var (
	verbose *bool
)

// OptionalLogger only logs when the verbose variable is true.
type OptionalLogger struct{}

// New creates a new VerboseLogger using the referenced bool to set verbosity.
func (l *OptionalLogger) New(v *bool) {
	verbose = v
}

// IsVerbose specifies if the logger should log human friendly output. This is nice for a program that has either human friendly verbose output or terse machine friendly return codes.
func (l *OptionalLogger) IsVerbose() bool {
	return *verbose
}

// Println prints a line if verbose is true.
func (l *OptionalLogger) Println(v ...interface{}) {
	if *verbose {
		log.Println(v)
	}
}

// Printf prints if verbose is true.
func (l *OptionalLogger) Printf(s string, v ...interface{}) {
	if *verbose {
		log.Printf(s, v)
	}
}
