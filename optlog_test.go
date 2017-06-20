// optlog is just a logger with 2 levels.
package ght_test

import (
	"testing"

	"github.com/ramjac/ght"
)

func TestNew(t *testing.T) {
	var vlog ght.OptionalLogger

	vlog.New(&[]bool{true}[0])

	if !vlog.IsVerbose() {
		t.Errorf("Logger was not verbose but should be")
	}
}
