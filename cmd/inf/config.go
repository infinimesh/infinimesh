package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/manifoldco/promptui"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

// TODO: manage multiple installations
type Config struct {
	CurrentContext   string     `yaml:"current-context,omitempty"`
	Contexts         []*Context `yaml:"contexts,omitempty"`
	DefaultNamespace string     `yaml:"defaultNamespace,omitempty"`
}

type Context struct {
	Name   string `yaml:"name"`
	Server string `yaml:"server"`
	TLS    bool   `yaml:"tls"`
	Token  string `yaml:"token,omitempty"`
}

var (
	apiserverFlag string
	tlsFlag       bool
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetContextCmd)
	configCmd.AddCommand(configSelectContext)
	configSetContextCmd.Flags().StringVar(&apiserverFlag, "apiserver", "grpc.infinimesh.io:443", "Infinimesh APIServer. Defaults to grpc.infinimesh.io:443")
	configSetContextCmd.Flags().BoolVar(&tlsFlag, "tls", true, "Enable or disable TLS. Defaults to true.")
}

var configSelectContext = &cobra.Command{
	Use:   "select-context",
	Short: "Interactively select a context",
	Run: func(cmd *cobra.Command, args []string) {
		var contextNames []string
		for _, context := range config.Contexts {
			contextNames = append(contextNames, context.Name)
		}
		p := promptui.Select{
			Label: "Select cluster",
			Items: contextNames,
		}

		_, selected, err := p.Run()
		if err != nil {
			os.Exit(0)
		}

		// How to have selection on currently selected cluster?

		config.CurrentContext = selected
		err = config.Write()
		if err != nil {
			panic(err)
		}
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure infinimesh CLI",
}

var configSetContextCmd = &cobra.Command{
	Use:   "set-context",
	Short: "Set fields in a context",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		newCtx := &Context{
			Name:   args[0],
			Server: apiserverFlag,
			TLS:    tlsFlag,
		}
		var replaced bool
		for i, ctx := range config.Contexts {
			if ctx.Name == args[0] {
				config.Contexts[i] = newCtx
				replaced = true
				fmt.Printf("Context %v updated.\n", args[0])
				break
			}
		}

		if !replaced {
			config.Contexts = append(config.Contexts, newCtx)
			fmt.Printf("Context %v created.\n", args[0])
		}

		err := config.Write()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write config file: %v.\n", err)
		}
	},
}

func (c *Config) GetCurrentContext() (*Context, error) {
	for _, ctx := range c.Contexts {
		if ctx.Name == c.CurrentContext {
			return ctx, nil
		}
	}

	return nil, errors.New("There is no context right now")
}

func (c *Config) Write() error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(home, ".inf")
	_ = os.MkdirAll(configDir, 0755)
	configPath := filepath.Join(configDir, "config")

	file, err := os.OpenFile(configPath, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	encoder := yaml.NewEncoder(file)
	return encoder.Encode(&c)
}

func ReadConfig() (c *Config, err error) {
	file, err := os.OpenFile(getDefaultConfigPath(), os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func getDefaultConfigPath() string {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(home, ".inf", "config")
}
