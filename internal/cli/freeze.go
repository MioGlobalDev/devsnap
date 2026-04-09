package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"devenv-snapshot/internal/detect"
	"devenv-snapshot/internal/detectors"
	"devenv-snapshot/internal/engine"
	"devenv-snapshot/internal/paths"
	"devenv-snapshot/internal/snapshot"
	"devenv-snapshot/internal/toolchain"

	"github.com/spf13/cobra"
)

type freezeFlags struct {
	root string
}

func NewFreezeCmd() *cobra.Command {
	var ff freezeFlags

	cmd := &cobra.Command{
		Use:     "freeze",
		Aliases: []string{"init"},
		Short:   "Generate .devenv snapshot from this project",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := filepath.Abs(ff.root)
			if err != nil {
				return err
			}

			dets, steps, err := engine.DetectAll(root, []engine.Detector{
				detectors.PythonDetector{},
				detectors.NodeDetector{},
			})
			if err != nil {
				return err
			}
			if len(dets) == 0 {
				return fmt.Errorf("no supported project detected (need Node lockfile or Python requirements/poetry)")
			}

			if err := os.MkdirAll(paths.Dir(root), 0o755); err != nil {
				return err
			}

			sd := snapshot.Detections{}
			for _, d := range dets {
				switch d.Kind {
				case "node":
					nd := d.Data.(detect.NodeDetection)
					sd.Node = &snapshot.Node{
						PackageManager: nd.PackageManager,
						LockFile:       filepath.ToSlash(nd.LockFile),
					}
				case "python":
					pd := d.Data.(detect.PythonDetection)
					sd.Python = &snapshot.Python{
						Manager: pd.Manager,
						File:    filepath.ToSlash(pd.File),
					}
				}
			}

			var ss []snapshot.Step
			for _, st := range steps {
				ss = append(ss, snapshot.Step{Run: st.Run})
			}

			s := snapshot.New(root, sd, ss)
			if sd.Node != nil {
				if v, ok, err := toolchain.NodeVersion(cmd.Context()); err == nil && ok {
					s.Toolchains.Node = v
				}
			}
			if sd.Python != nil {
				if v, ok, err := toolchain.PythonVersion(cmd.Context()); err == nil && ok {
					s.Toolchains.Python = v
				}
			}

			if err := snapshot.WriteFile(paths.SnapshotPath(root), s); err != nil {
				return err
			}

			if err := writeRestoreScripts(root); err != nil {
				return err
			}

			fmt.Fprintf(os.Stdout, "wrote %s\n", paths.SnapshotPath(root))
			return nil
		},
	}

	cmd.Flags().StringVar(&ff.root, "root", ".", "project root directory")
	return cmd
}

func writeRestoreScripts(root string) error {
	ps1 := `Param(
  [string]$ProjectRoot = "."
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

Push-Location $ProjectRoot
try {
  if (Test-Path .\devsnap.exe) {
    .\devsnap.exe restore --root .
  } elseif (Test-Path .\devenv.exe) {
    .\devenv.exe restore --root .
  } elseif (Get-Command devsnap -ErrorAction SilentlyContinue) {
    devsnap restore --root .
  } elseif (Get-Command devenv -ErrorAction SilentlyContinue) {
    devenv restore --root .
  } else {
    throw "devsnap not found. Build or download devsnap, then run: devsnap restore --root ."
  }
} finally {
  Pop-Location
}
`

	sh := `#!/usr/bin/env sh
set -eu

PROJECT_ROOT="${1:-.}"
cd "$PROJECT_ROOT"

if [ -x "./devsnap" ]; then
  ./devsnap restore --root .
elif [ -x "./devenv" ]; then
  ./devenv restore --root .
elif command -v devsnap >/dev/null 2>&1; then
  devsnap restore --root .
elif command -v devenv >/dev/null 2>&1; then
  devenv restore --root .
else
  echo "devsnap not found. Build or download devsnap, then run: devsnap restore --root ." >&2
  exit 1
fi
`

	if err := os.WriteFile(paths.RestorePS1Path(root), []byte(ps1), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(paths.RestoreSHPath(root), []byte(sh), 0o755); err != nil {
		return err
	}
	return nil
}

