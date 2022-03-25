/*
Copyright © 2022 Infinite Devices GmbH, Nikita Ivanovski info@slnt-opp.xyz

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "inf",
	Short: "infinimesh Platform CLI",
	RunE: contextCmd.RunE,
}

var VERSION string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	VERSION = version
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.default.infinimesh.yaml)")
	rootCmd.PersistentFlags().Bool("json", false, "Print output as json")
	rootCmd.PersistentFlags().Bool("verbose", false, "Print additional info related to the CLI itself")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".inf" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".default.infinimesh")

		cfgFile = fmt.Sprintf("%s/.default.infinimesh.yaml", home)
	}

	viper.AutomaticEnv() // read in environment variables that match

	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		if _, err := os.Create(cfgFile); err != nil { // perm 0666
			fmt.Fprintln(os.Stderr, "Can't create default config file")
			panic(err)
		}
	}

	verbose, _ := rootCmd.Flags().GetBool("verbose")
	// If a config file is found, read it in.
	err := viper.ReadInConfig();
	if err == nil && verbose {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func printJsonResponse(data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))
	return nil
}