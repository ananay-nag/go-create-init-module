package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config struct to hold YAML config values
type Config struct {
	PreSet string `yaml:"pre-set"`
}

// LoadConfig searches for mod-name.yaml in parent directories
func LoadConfig() (Config, string, error) {
	var config Config

	// Start search from current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return config, "", err
	}

	for {
		configPath := filepath.Join(currentDir, "mod-name.yaml")
		if _, err := os.Stat(configPath); err == nil {
			file, err := os.ReadFile(configPath)
			if err != nil {
				return config, "", err
			}
			err = yaml.Unmarshal(file, &config)
			return config, currentDir, err
		}

		// Move up one directory
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break // Reached filesystem root, stop search
		}
		currentDir = parentDir
	}

	return config, "", fmt.Errorf("mod-name.yaml not found")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: set-mod <module-name> or set-mod -c to init in current dir")
		os.Exit(1)
	}

	config, projectRoot, err := LoadConfig()
	if err != nil {
		log.Fatalf("Error loading mod-name.yaml: %v", err)
	}

	if config.PreSet == "" {
		log.Fatalf("Config 'pre-set' value is missing.")
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}

	// Compute relative path from project root to cwd
	relativePath, err := filepath.Rel(projectRoot, cwd)
	if err != nil {
		log.Fatalf("Error computing relative path: %v", err)
	}

	// Clean up relative path
	relativePath = strings.TrimPrefix(relativePath, "./")
	relativePath = strings.TrimPrefix(relativePath, ".")

	// Determine module path
	var modulePath string

	if os.Args[1] == "-c" {
		// Use the current directory name directly
		currentDirName := filepath.Base(cwd)
    		modulePath = fmt.Sprintf("%s/%s", config.PreSet, currentDirName)
	} else {
		moduleName := os.Args[1]
		modulePath = fmt.Sprintf("%s/%s/%s", config.PreSet, relativePath, moduleName)

		// Ensure module directory exists (only for named modules)
		if err := os.MkdirAll(filepath.Join(cwd, moduleName), os.ModePerm); err != nil {
			log.Fatalf("Error creating module directory: %v", err)
		}

		// Change working directory to the new module
		cwd = filepath.Join(cwd, moduleName)
	}

	// Remove unnecessary slashes
	modulePath = strings.TrimRight(modulePath, "/")
	modulePath = strings.ReplaceAll(modulePath, "//", "/")

	// Ensure no leading or trailing slashes
	modulePath = strings.TrimPrefix(modulePath, "/")
	modulePath = strings.TrimSuffix(modulePath, "/")

	// Run `go mod init`
	cmd := exec.Command("go", "mod", "init", modulePath)
	cmd.Dir = cwd // Use the correct working directory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Initializing Go module at %s\n", modulePath)
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error running 'go mod init': %v", err)
	}

	fmt.Println("Go module initialized successfully!")
}
