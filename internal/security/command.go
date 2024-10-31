// internal/security/command.go
package security

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// Command types for static analysis
type Command struct {
	Path string
	Name string
}

// Define constants for supported commands
const (
	CmdNanoDNS = "nanodns"
	CmdCat     = "cat" // Used for testing
)

var (
	// Error definitions
	ErrInvalidCommand = errors.New("invalid command")
	ErrInvalidArgs    = errors.New("invalid command arguments")

	// Pre-defined set of allowed commands (immutable after init)
	allowedCommands map[string]bool

	// Pre-defined set of shell metacharacters (immutable)
	shellMetacharacters = []string{";", "|", "&", ">", "<", "`", "$", "(", ")", "\n", "\r"}
)

func init() {
	// Initialize allowed commands
	allowedCommands = map[string]bool{
		CmdNanoDNS: true,
		CmdCat:     true, // Only included when testing
	}
}

// NewCommand creates a new Command instance
func NewCommand(name string) Command {
	return Command{
		Name: name,
	}
}

// SecureCommand creates a validated and sanitized exec.Cmd
func SecureCommand(name string, args ...string) (*exec.Cmd, error) {
	cmd := NewCommand(name)
	return cmd.CreateExecCommand(args...)
}

// CreateExecCommand creates an exec.Cmd with security validations
func (c Command) CreateExecCommand(args ...string) (*exec.Cmd, error) {
	// 1. Check for shell metacharacters in command name
	if containsShellMetacharacters(c.Name) {
		return nil, fmt.Errorf("%w: command contains invalid characters", ErrInvalidCommand)
	}

	// 2. Validate command name against allowed list
	baseName := filepath.Base(c.Name)
	if !allowedCommands[baseName] {
		return nil, fmt.Errorf("%w: %s", ErrInvalidCommand, c.Name)
	}

	// 3. Find command in PATH
	fullPath, err := exec.LookPath(c.Name)
	if err != nil {
		return nil, fmt.Errorf("command not found: %s", err)
	}

	// 4. Validate arguments
	if err := validateArgs(args); err != nil {
		return nil, err
	}

	// 5. Create command with validated inputs
	return exec.Command(fullPath, args...), nil
}

// validateArgs checks command arguments for security issues
func validateArgs(args []string) error {
	for _, arg := range args {
		if containsShellMetacharacters(arg) {
			return fmt.Errorf("%w: contains shell metacharacters", ErrInvalidArgs)
		}
	}
	return nil
}

// containsShellMetacharacters checks for dangerous shell characters
func containsShellMetacharacters(s string) bool {
	for _, char := range shellMetacharacters {
		if strings.Contains(s, char) {
			return true
		}
	}
	return false
}

// For testing purposes only
type testingConfig struct{}

var testConfig testingConfig

// GetTestConfig returns the test configuration helper
func GetTestConfig() testingConfig {
	return testConfig
}
