// internal/security/command_test.go
package security

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func setupTestEnvironment(t *testing.T) (string, func()) {
	t.Helper()

	// Create temporary test directory
	tmpDir := t.TempDir()
	binDir := filepath.Join(tmpDir, "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		t.Fatalf("failed to create bin dir: %v", err)
	}

	// Create mock executables
	mockFiles := []string{"nanodns", "cat"}
	for _, name := range mockFiles {
		var scriptContent string
		if runtime.GOOS == "windows" {
			scriptContent = `@echo off
echo test content`
		} else {
			scriptContent = `#!/bin/sh
echo test content`
		}

		ext := ""
		if runtime.GOOS == "windows" {
			ext = ".bat"
		}

		path := filepath.Join(binDir, name+ext)
		if err := os.WriteFile(path, []byte(scriptContent), 0755); err != nil {
			t.Fatalf("failed to create mock command %s: %v", name, err)
		}
	}

	// Save original PATH and update it
	origPath := os.Getenv("PATH")
	newPath := binDir + string(os.PathListSeparator) + origPath
	os.Setenv("PATH", newPath)

	return binDir, func() {
		os.Setenv("PATH", origPath)
	}
}

func TestSecureCommand(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	tests := []struct {
		name        string
		cmdName     string
		args        []string
		expectError bool
		errorType   error
	}{
		{
			name:        "Valid command without args",
			cmdName:     "nanodns",
			args:        []string{},
			expectError: false,
		},
		{
			name:        "Command with shell injection",
			cmdName:     "nanodns;rm -rf /",
			args:        []string{},
			expectError: true,
			errorType:   ErrInvalidCommand,
		},
		{
			name:        "Invalid command",
			cmdName:     "malicious",
			args:        []string{},
			expectError: true,
			errorType:   ErrInvalidCommand,
		},
		{
			name:        "Command with shell injection in args",
			cmdName:     "nanodns",
			args:        []string{"-config", "config.yaml; rm -rf /"},
			expectError: true,
			errorType:   ErrInvalidArgs,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd, err := SecureCommand(tc.cmdName, tc.args...)

			if tc.expectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tc.errorType != nil && err != nil && !errors.Is(err, tc.errorType) {
				t.Errorf("expected error type %v but got %v", tc.errorType, err)
			}
			if !tc.expectError && cmd != nil {
				if err := cmd.Run(); err != nil {
					t.Errorf("command execution failed: %v", err)
				}
			}
		})
	}
}
