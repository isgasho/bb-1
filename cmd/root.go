package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/craftamap/bb/cmd/commands/api"
	"github.com/craftamap/bb/cmd/commands/auth"
	"github.com/craftamap/bb/cmd/commands/downloads"
	"github.com/craftamap/bb/cmd/commands/pr"
	"github.com/craftamap/bb/cmd/options"
	"github.com/kirsle/configdir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:     "bb",
		Short:   "Bitbucket.org CLI",
		Long:    "Work seamlessly with Bitbucket.org from the command line.",
		Example: `$ bb pr list`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			viper.Unmarshal(&globalOpts)
		},
	}

	cfgFile    string
	globalOpts = options.GlobalOptions{}

	username string
	password string
	repoOrga string
	repoSlug string
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/bb)")
	rootCmd.PersistentFlags().StringVar(&username, "username", "", "username")
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "app password")

	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))

	pr.Add(rootCmd, &globalOpts)
	api.Add(rootCmd, &globalOpts)
	downloads.Add(rootCmd, &globalOpts)
	auth.Add(rootCmd, &globalOpts)
}

func initConfig() {
	if cfgFile == "" {
		configDir := configdir.LocalConfig("bb")
		err := configdir.MakePath(configDir)
		if err != nil {
			panic(err)
		}
		cfgFile = filepath.Join(configDir, "configuration.toml")
		if _, err = os.Stat(cfgFile); os.IsNotExist(err) {
			fh, err := os.Create(cfgFile)
			if err != nil {
				panic(err)
			}
			defer fh.Close()
		}
	}

	viper.SetConfigFile(cfgFile)

	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}
