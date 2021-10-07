package cmd

import (
	"github.com/spf13/cobra"

	"github.com/thibaultmg/clingua/internal/config"
)

var (
	ClinguaVersion = "development"

	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "clingua",
		Short: "A CLI application for creating and learning vocabulary cards.",
		Long: `Clingua is a CLI application integrating dictionnaries, translators
and Speech APIs to create easily vocabulary cards.
Those cards can then be learned through generated quizzes based on those cards.`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Automatically adds --version command
	rootCmd.Version = ClinguaVersion
	cobra.OnInitialize(func() {
		config.InitConfig(cfgFile)
	})

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $PWD/.clingua.yaml)")
}
