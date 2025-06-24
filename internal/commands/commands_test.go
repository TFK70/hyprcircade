package commands

import (
	"testing"

	"github.com/tfk70/hyprcircade/internal/logging"
)

func TestExecuteCommand(t *testing.T) {
	logging.SetupLogger()

	cmd := "echo 10"
	err := ExecuteCommand(cmd)
	if err != nil {
		t.Fatalf("Error during cmd execution\nCmd: %s\nError: %s", cmd, err.Error())
	}

	cmd = "this-command-does-not-exist-at-all it-really-doesnt"
	err = ExecuteCommand(cmd)
	if err == nil {
		t.Fatalf("Error should be thrown while executing command that does not exist")
	}
}
