// optlog is just a logger with 2 levels.

package ght

import (
	"github.com/fatih/color"
)

var (
	verbose *bool
	logger  *color.Color
)

// VerboseLogger only logs when the verbose variable is true.
type VerboseLogger struct{}

// New creates a new VerboseLogger using the referenced bool to set verbosity.
func (l *VerboseLogger) New(v *bool) {
	verbose = v
	logger = color.New()
}

// SetColor sets the color of the logger's output.
func (l *VerboseLogger) SetColor(p ...color.Attribute) {
	logger = color.Set(p...)
}

// IsVerbose specifies if the logger should log human friendly output. This is nice for a program that has either human friendly verbose output or terse machine friendly return codes.
func (l *VerboseLogger) IsVerbose() bool {
	return *verbose
}

// Println prints a line if verbose is true.
func (l *VerboseLogger) Println(v ...interface{}) {
	if *verbose {
		logger.Println(v)
	}
}

// Printf prints if verbose is true.
func (l *VerboseLogger) Printf(s string, v ...interface{}) {
	if *verbose {
		logger.Printf(s, v)
	}
}
