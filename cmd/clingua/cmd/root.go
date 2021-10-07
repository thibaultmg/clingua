package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/thibaultmg/clingua/internal/config"
)

var (
	ClinguaVersion = "development"

	cfgFile  string
	logLevel string

	rootCmd = &cobra.Command{
		Use:   "clingua",
		Short: "A CLI application for creating and learning vocabulary cards.",
		Long: `Clingua is a CLI application integrating dictionnaries, translators
and Speech APIs to create easily vocabulary cards.
Those cards can then be learned through generated quizzes based on those cards.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setLogLevel()
		},
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
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "warn", "log level")
	viper.BindPFlag("logLevel", rootCmd.PersistentFlags().Lookup("log-level"))

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

}

func setLogLevel() {
	zlLevel, err := zerolog.ParseLevel(viper.GetString("logLevel"))
	if err != nil {
		log.Fatal().Err(err).Msg("Invalid log level")
	}

	log.Debug().Msgf("setting log level to %v", zlLevel)

	zerolog.SetGlobalLevel(zlLevel)
}
