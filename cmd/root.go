package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

const (
	CmdName           = "dcfs"
	DefaultConfigFile = CmdName + ".yaml"
)

var (
	cfg          = NewConfig()
	cfgFile      string
	cfgReadError error
)

var rootCmd = &cobra.Command{
	Use: CmdName,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// root level flags
	pflags := rootCmd.PersistentFlags()
	pflags.StringVarP(&cfgFile, "config-file", "f", DefaultConfigFile, "config file (YAML format)")

	// load config-file before cobra commands
	cobra.OnInitialize(func() {
		if cfgFile != "" {
			if err := cfg.ReadInFile(cfgFile); err != nil {
				cfgReadError = err
				log.Println(err)
			}
		}
	})
}
