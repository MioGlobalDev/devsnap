package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"devenv-snapshot/internal/paths"
	"devenv-snapshot/internal/snapshot"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type printFlags struct {
	root string
}

func NewPrintCmd() *cobra.Command {
	var pf printFlags

	cmd := &cobra.Command{
		Use:   "print",
		Short: "Print the parsed .devenv snapshot",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := filepath.Abs(pf.root)
			if err != nil {
				return err
			}

			s, err := snapshot.ReadFile(paths.SnapshotPath(root))
			if err != nil {
				return err
			}

			b, err := yaml.Marshal(&s)
			if err != nil {
				return err
			}
			fmt.Fprint(os.Stdout, string(b))
			return nil
		},
	}

	cmd.Flags().StringVar(&pf.root, "root", ".", "project root directory")
	return cmd
}

