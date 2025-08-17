package main

import (
	"bufio"
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
			// File found, read and unmarshal it
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

	var config Config
	var projectRoot string
	var err error

	// Try to load the config from an existing file
	config, projectRoot, err = LoadConfig()
	if err != nil {
		// If the file is not found, create a new one in the current directory
		if strings.Contains(err.Error(), "mod-name.yaml not found") {
			cwd, createErr := os.Getwd()
			if createErr != nil {
				log.Fatalf("Error getting current directory to create config: %v", createErr)
			}
			configPath := filepath.Join(cwd, "mod-name.yaml")
			defaultPreSet := "github.com/your-username"

			fmt.Printf("mod-name.yaml not found. Creating a new one at %s\n", configPath)
			fmt.Printf("Please enter a value for 'pre-set' (default: %s): ", defaultPreSet)

			// Read user input
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			userInput := strings.TrimSpace(input)

			// Use user input if it's not empty, otherwise use the default
			preSetValue := defaultPreSet
			if userInput != "" {
				preSetValue = userInput
			}

			defaultContent := []byte(fmt.Sprintf(`pre-set: "%s"`, preSetValue))

			if writeErr := os.WriteFile(configPath, defaultContent, 0644); writeErr != nil {
				log.Fatalf("Error creating mod-name.yaml: %v", writeErr)
			}

			// Now that the file is created, try loading the config again
			config, projectRoot, err = LoadConfig()
			if err != nil {
				// This should not happen, but we handle it just in case
				log.Fatalf("Error loading newly created mod-name.yaml: %v", err)
			}
		} else {
			// Handle other types of errors from LoadConfig
			log.Fatalf("Error loading mod-name.yaml: %v", err)
		}
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
