package data

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/davidroman0O/gogog/types"
)

func getPath() (appDataPath string) {
	switch runtime.GOOS {
	case "windows":
		appDataPath = os.Getenv("APPDATA") + "/gogog"
	case "darwin": // macOS
		home, _ := os.UserHomeDir()
		appDataPath = home + "/Library/Application Support" + "/gogog"
	default: // Assume Linux/Unix
		home, _ := os.UserHomeDir()
		appDataPath = home + "/.local/share" + "/gogog"
	}
	return appDataPath
}

func getDataPath[T types.GogStates](client bool) string {
	name := strings.ToLower(reflect.TypeOf(*new(T)).Name())
	if client {
		return fmt.Sprintf("%v/client/%v.json", getPath(), name)
	}
	return fmt.Sprintf("%v/agent/%v.json", getPath(), name)
}

func fileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if err == nil {
		return !info.IsDir()
	}
	if os.IsNotExist(err) {
		return false
	}
	return false // File may exist but there's an error accessing it
}

func writeFile(filePath string, data []byte) error {
	// Check if the file exists
	if _, err := os.Stat(filePath); err == nil {
		// File exists, handle accordingly
		// For example, you can return an error or prompt the user
		return os.WriteFile(filePath, data, 0644)
		// return os.ErrExist
	} else if !os.IsNotExist(err) {
		// Some other error occurred while checking the file
		return err
	}

	// Create directories recursively if they don't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Write data to file (create if not exists)
	return os.WriteFile(filePath, data, 0644) // 0644 permissions
}

var asClient bool = true

func SetRuntime(client bool) {
	asClient = client
}

func IsClient() bool {
	return asClient
}

func Has[T types.GogStates]() bool {
	return fileExists(getDataPath[T](IsClient()))
}

func Load[T types.GogStates]() (state *T, err error) {
	var appDataPath string = getDataPath[T](IsClient())
	// Read data from file
	data, err := os.ReadFile(appDataPath)
	if err != nil {
		log.Fatalf("Failed to read from file: %v", err)
		return nil, err
	}
	return state, json.Unmarshal(data, state)
}

func Delete[T types.GogStates]() error {
	filePath := getDataPath[T](IsClient())
	err := os.Remove(filePath)
	if err != nil {
		log.Fatalf("Failed to delete file: %v", err)
		return err
	}
	return nil
}

func Save[T types.GogStates](state *T) error {
	var err error
	var bytesState []byte
	var appDataPath string = getDataPath[T](IsClient())
	fmt.Println("App Data Path:", appDataPath)
	if bytesState, err = json.Marshal(*state); err != nil {
		return err
	}
	if err = writeFile(appDataPath, bytesState); err != nil {
		if err == os.ErrExist {
			fmt.Printf("File already exists: %s\n", appDataPath)
		} else {
			fmt.Printf("Failed to write to file: %v\n", err)
		}
		return err
	}
	return nil
}
