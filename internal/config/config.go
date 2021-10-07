package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
	"golang.org/x/text/language"
)

func InitConfig(cfgFile string) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find working directory.
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Error getting current working dir: %v", err)
		}

		// Search config in working directory.
		viper.AddConfigPath(pwd)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".clingua")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func GetOxfordRepo() (url, appID, appKey string) {
	url = viper.GetString("oxford.url")
	appID = viper.GetString("oxford.appID")
	appKey = viper.GetString("oxford.appKey")

	return
}

func GetLanguages() (from, to language.Tag) {
	from = language.MustParse(viper.GetString("fromLanguage"))
	to = language.MustParse(viper.GetString("toLanguage"))

	return
}
