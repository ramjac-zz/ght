// optlog is just a logger with 2 levels.

package ght

import "log"

var (
	verbose *bool
)

// VerboseLogger only logs when the verbose variable is true.
type VerboseLogger struct{}

// New creates a new VerboseLogger using the referenced bool to set verbosity.
func (l *VerboseLogger) New(v *bool) {
	verbose = v
}

func (l *VerboseLogger) IsVerbose() bool {
	return *verbose
}

// Println prints a line if verbose is true.
func (l *VerboseLogger) Println(v ...interface{}) {
	if *verbose {
		log.Println(v)
	}
}

// Printf prints a line if verbose is true.
func (l *VerboseLogger) Printf(s string, v ...interface{}) {
	if *verbose {
		log.Printf(s, v)
	}
}
