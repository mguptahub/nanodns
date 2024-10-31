// internal/security/command.go
package security

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	// Allowed commands with default values
	allowedCommands = map[string]bool{
		"nanodns": true,
		"cat":     true, // Added for testing
	}

	ErrInvalidCommand = errors.New("invalid command")
	ErrInvalidPath    = errors.New("invalid command path")
	ErrInvalidArgs    = errors.New("invalid command arguments")
)

// SecureCommand creates a validated and sanitized exec.Cmd
func SecureCommand(name string, args ...string) (*exec.Cmd, error) {
	// 1. Validate command name
	baseName := filepath.Base(name)
	if !allowedCommands[baseName] {
		return nil, fmt.Errorf("%w: %s", ErrInvalidCommand, name)
	}

	// 2. Find command in PATH
	fullPath, err := exec.LookPath(name)
	if err != nil {
		return nil, fmt.Errorf("command not found: %s", err)
	}

	// 3. Validate command path
	if err := ValidateCommandPath(fullPath); err != nil {
		return nil, err
	}

	// 4. Validate arguments
	if err := validateArgs(args); err != nil {
		return nil, err
	}

	// 5. Create command with validated inputs
	cmd := exec.Command(fullPath, args...)
	return cmd, nil
}

// ValidateCommandPath checks if a command path is allowed
func ValidateCommandPath(path string) error {
	// Clean and resolve the path
	cleanPath := filepath.Clean(path)

	// Check if path is absolute
	if !filepath.IsAbs(cleanPath) {
		return fmt.Errorf("%w: relative paths not allowed", ErrInvalidPath)
	}

	// Check for shell metacharacters in path
	if containsShellMetacharacters(cleanPath) {
		return fmt.Errorf("%w: path contains invalid characters", ErrInvalidPath)
	}

	return fmt.Errorf("%w: %s", ErrInvalidPath, path)
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
	metachars := []string{";", "|", "&", ">", "<", "`", "$", "(", ")", "\n", "\r"}
	for _, char := range metachars {
		if strings.Contains(s, char) {
			return true
		}
	}
	return false
}
