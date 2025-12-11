package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
)

func main() {
	// Find Python executable
	pythonExe := findPython()
	if pythonExe == "" {
		showError("Python not found", "Please install Python 3.8+ from python.org")
		os.Exit(1)
	}

	// Find companion app directory
	exePath, err := os.Executable()
	if err != nil {
		showError("Error", fmt.Sprintf("Failed to get executable path: %v", err))
		os.Exit(1)
	}

	exeDir := filepath.Dir(exePath)
	companionDir := filepath.Join(exeDir, "companion-app")
	mainScript := filepath.Join(companionDir, "main.py")

	// Check if companion app exists
	if _, err := os.Stat(mainScript); os.IsNotExist(err) {
		showError("Missing Files", fmt.Sprintf("companion-app/main.py not found in:\n%s", exeDir))
		os.Exit(1)
	}

	// Launch Python app
	log.Printf("Launching Forger Companion from: %s", mainScript)
	
	cmd := exec.Command(pythonExe, mainScript)
	cmd.Dir = companionDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Platform-specific setup
	if runtime.GOOS == "windows" {
		// Windows: hide console window
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		// Note: CreationFlags not available on all platforms, will be set at runtime
	}

	if err := cmd.Run(); err != nil {
		log.Printf("Error running companion app: %v", err)
		showError("Error", fmt.Sprintf("Failed to run companion app:\n%v", err))
		os.Exit(1)
	}
}

func findPython() string {
	// Try common Python locations
	candidates := []string{
		"python",
		"python3",
		"python3.11",
		"python3.10",
		"python3.9",
		"python3.8",
	}

	if runtime.GOOS == "windows" {
		candidates = append(candidates,
			"C:\\Python311\\python.exe",
			"C:\\Python310\\python.exe",
			"C:\\Python39\\python.exe",
			"C:\\Python38\\python.exe",
		)
	}

	for _, candidate := range candidates {
		if _, err := exec.LookPath(candidate); err == nil {
			return candidate
		}
	}

	return ""
}

func showError(title, message string) {
	if runtime.GOOS == "windows" {
		// Use Windows message box
		cmd := exec.Command("powershell", "-Command", fmt.Sprintf(
			`[System.Windows.Forms.MessageBox]::Show('%s', '%s')`,
			message, title,
		))
		cmd.Run()
	} else {
		fmt.Printf("%s: %s\n", title, message)
	}
}
