package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// ExecuteCommand executes a shell command and returns its output
func ExecuteCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("command failed: %w\nstderr: %s", err, stderr.String())
	}
	
	return stdout.String(), nil
}

// CommandExists checks if a command exists in the system
func CommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// ShellCommand executes a shell command with sh -c
func ShellCommand(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("shell command failed: %w\nstderr: %s", err, stderr.String())
	}
	
	return stdout.String(), nil
}

// GetCommandOutput executes a command and ignores errors
func GetCommandOutput(name string, args ...string) string {
	output, _ := ExecuteCommand(name, args...)
	return strings.TrimSpace(output)
}

// CheckRootPrivileges checks if the application is running with root privileges
func CheckRootPrivileges() bool {
	output := GetCommandOutput("id", "-u")
	return output == "0"
}
