package config

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// mustReadConfigFile checks if a file exists in one of the configuration
// directories and returns the content. If no file is found, we exit.
func mustReadConfigFile(filename string, locations []string) ([]byte, string) {
	for _, dir := range locations {
		if strings.HasPrefix(dir, "~") {
			homeDir := os.Getenv("HOME")
			dir = filepath.Join(homeDir, dir[1:])
		}
		filep := filepath.Join(dir, filename)
		log.WithField("filepath", filep).Debug("Looking for config file")
		if fileExists(filep) {
			return mustReadFile(filep), dir
		}
	}
	errMsg := "Could not find config file"
	if len(locations) > 1 {
		errMsg += " in any of the possible directories"
	}
	log.WithField("filepath", filename).Fatal(errMsg)
	return nil, ""
}

// fileExists checks if a given file exists.
func fileExists(path string) bool {
	if _, err := os.Stat(evalSymlink(path)); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		log.WithError(err).Error()
		return false
	}
}

func evalSymlink(path string) string {
	if path == "" {
		return path
	}
	if path[0] == '~' {
		path = os.Getenv("HOME") + path[1:]
	}
	evalPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return path
	}
	return evalPath
}

// readFile reads a given file and returns the content.
func readFile(filename string) ([]byte, error) {
	filename = evalSymlink(filename)
	log.WithField("filepath", filename).Trace("Reading file...")
	return os.ReadFile(filename)
}

// MustReadFile reads a given file and returns the content. If an error occurs the application terminates.
func mustReadFile(filename string) []byte {
	file, err := readFile(filename)
	if err != nil {
		log.WithError(err).Fatal("Error reading config file")
	}
	log.WithField("filepath", filename).Info("Read config file")
	return file
}
