package snapshot

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"gopkg.in/yaml.v3"
)

type Snapshot struct {
	SchemaVersion string    `yaml:"schemaVersion"`
	CreatedAt     time.Time `yaml:"createdAt"`
	Platform      Platform  `yaml:"platform"`
	Project       Project   `yaml:"project"`
	Toolchains    Toolchains `yaml:"toolchains,omitempty"`
	Detections    Detections `yaml:"detections"`
	Steps         []Step     `yaml:"steps"`
}

type Platform struct {
	OS   string `yaml:"os"`
	Arch string `yaml:"arch"`
}

type Project struct {
	Root string `yaml:"root"`
}

type Toolchains struct {
	Node   string `yaml:"node,omitempty"`
	Python string `yaml:"python,omitempty"`
}

type Detections struct {
	Node   *Node   `yaml:"node,omitempty"`
	Python *Python `yaml:"python,omitempty"`
}

type Node struct {
	PackageManager string `yaml:"packageManager"`
	LockFile       string `yaml:"lockFile"`
}

type Python struct {
	Manager string `yaml:"manager"` // pip | poetry
	File    string `yaml:"file"`    // requirements.txt | poetry.lock
}

type Step struct {
	Run string `yaml:"run"`
}

func New(root string, detections Detections, steps []Step) Snapshot {
	return Snapshot{
		SchemaVersion: "0.3",
		CreatedAt:     time.Now().UTC(),
		Platform: Platform{
			OS:   runtime.GOOS,
			Arch: runtime.GOARCH,
		},
		Project: Project{
			Root: filepath.ToSlash(root),
		},
		Detections: detections,
		Steps:      steps,
	}
}

func WriteFile(path string, s Snapshot) error {
	b, err := yaml.Marshal(&s)
	if err != nil {
		return err
	}
	b = append([]byte("# DevEnv Snapshot (v0.3)\n"), b...)
	return os.WriteFile(path, b, 0o644)
}

func ReadFile(path string) (Snapshot, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, err
	}
	var s Snapshot
	if err := yaml.Unmarshal(b, &s); err != nil {
		return Snapshot{}, err
	}
	return s, nil
}

