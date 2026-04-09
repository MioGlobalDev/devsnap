package detect

import (
	"os"
	"path/filepath"
)

type NodeDetection struct {
	PackageManager string
	LockFile       string
}

func DetectNode(root string) (NodeDetection, bool, error) {
	// priority: pnpm > npm > yarn
	candidates := []struct {
		pm   string
		file string
	}{
		{pm: "pnpm", file: "pnpm-lock.yaml"},
		{pm: "npm", file: "package-lock.json"},
		{pm: "yarn", file: "yarn.lock"},
	}

	for _, c := range candidates {
		p := filepath.Join(root, c.file)
		if _, err := os.Stat(p); err == nil {
			return NodeDetection{
				PackageManager: c.pm,
				LockFile:       c.file,
			}, true, nil
		}
	}

	return NodeDetection{}, false, nil
}

