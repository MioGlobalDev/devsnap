package detect

import (
	"os"
	"path/filepath"
)

type PythonDetection struct {
	Manager string // "poetry" | "pip"
	File    string // "poetry.lock" | "requirements.txt"
}

func DetectPython(root string) (PythonDetection, bool, error) {
	// priority: poetry > requirements
	poetryLock := filepath.Join(root, "poetry.lock")
	pyproject := filepath.Join(root, "pyproject.toml")
	if exists(poetryLock) && exists(pyproject) {
		return PythonDetection{Manager: "poetry", File: "poetry.lock"}, true, nil
	}

	req := filepath.Join(root, "requirements.txt")
	if exists(req) {
		return PythonDetection{Manager: "pip", File: "requirements.txt"}, true, nil
	}

	return PythonDetection{}, false, nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

