package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"devenv-snapshot/internal/paths"
	"devenv-snapshot/internal/runner"
	"devenv-snapshot/internal/snapshot"
	"devenv-snapshot/internal/toolchain"

	"github.com/spf13/cobra"
)

type restoreFlags struct {
	root string
}

func NewRestoreCmd() *cobra.Command {
	var rf restoreFlags

	cmd := &cobra.Command{
		Use:   "restore",
		Short: "Restore dependencies according to .devenv snapshot",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := filepath.Abs(rf.root)
			if err != nil {
				return err
			}

			s, err := snapshot.ReadFile(paths.SnapshotPath(root))
			if err != nil {
				return fmt.Errorf("read snapshot: %w", err)
			}

			if s.Detections.Node != nil {
				lockPath := filepath.Join(root, filepath.FromSlash(s.Detections.Node.LockFile))
				if _, err := os.Stat(lockPath); err != nil {
					return fmt.Errorf("lockfile missing: %s (run `devsnap init` again after generating a lockfile)", s.Detections.Node.LockFile)
				}
			}
			if s.Detections.Python != nil {
				fp := filepath.Join(root, filepath.FromSlash(s.Detections.Python.File))
				if _, err := os.Stat(fp); err != nil {
					return fmt.Errorf("python manifest missing: %s (run `devsnap init` again after generating files)", s.Detections.Python.File)
				}
			}

			if len(s.Steps) == 0 {
				return fmt.Errorf("snapshot has no steps to run")
			}

			// Level 2 (minimal): version hinting only. No auto-install.
			if s.Toolchains.Node != "" {
				if cur, ok, err := toolchain.NodeVersion(cmd.Context()); err == nil && ok && cur != "" && cur != s.Toolchains.Node {
					fmt.Fprintf(os.Stdout, "[devsnap] Node version mismatch\n  expected: %s\n  current:  %s\n", s.Toolchains.Node, cur)
				}
			}
			if s.Toolchains.Python != "" {
				if cur, ok, err := toolchain.PythonVersion(cmd.Context()); err == nil && ok && cur != "" && cur != s.Toolchains.Python {
					fmt.Fprintf(os.Stdout, "[devsnap] Python version mismatch\n  expected: %s\n  current:  %s\n", s.Toolchains.Python, cur)
				}
			}

			for _, st := range s.Steps {
				fmt.Fprintf(os.Stdout, "run: %s\n", st.Run)
				res, err := runner.RunShell(context.Background(), root, st.Run, 30*time.Minute)
				if res.Stdout != "" {
					fmt.Fprintln(os.Stdout, res.Stdout)
				}
				if err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&rf.root, "root", ".", "project root directory")
	return cmd
}

