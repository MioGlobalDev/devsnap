package cli

import (
	"fmt"
	"os"

	"devenv-snapshot/internal/version"

	"github.com/spf13/cobra"
)

type rootFlags struct {
	verbose bool
}

func NewRootCmd() *cobra.Command {
	var rf rootFlags

	cmd := &cobra.Command{
		Use:           "devsnap",
		Aliases:       []string{"devenv"},
		Short:         "DevEnv Snapshot: freeze and restore a reproducible project state",
		Version:       version.Version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.Flags().BoolP("version", "V", false, "print version")
	cmd.SetVersionTemplate("devsnap {{.Version}}\n")
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		v, _ := cmd.Flags().GetBool("version")
		if v {
			fmt.Fprintf(cmd.OutOrStdout(), "devsnap %s\n", version.Version)
			os.Exit(0)
		}
		return nil
	}

	cmd.PersistentFlags().BoolVarP(&rf.verbose, "verbose", "v", false, "verbose output")

	cmd.AddCommand(NewFreezeCmd())
	cmd.AddCommand(NewRestoreCmd())
	cmd.AddCommand(NewDoctorCmd())
	cmd.AddCommand(NewPrintCmd())

	return cmd
}

