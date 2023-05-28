package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/KazukiHayase/forge/codegen"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the project",
	Long:  `Initialize the project by creating the .forge directory and generating the configuration file inside it.`,
	RunE:  runInitCmd,
}

func runInitCmd(_ *cobra.Command, _ []string) error {
	if err := os.MkdirAll(codegen.RootDir, os.ModePerm); err != nil {
		return err
	}

	config := codegen.Config{
		Name: "sample",
		Prompts: []codegen.Prompt{
			{Name: "prompt-1", Message: "input prompt-1"},
			{Name: "prompt-2", Message: "input prompt-2"},
		},
		InOuts: []codegen.InOut{
			{Input: "input-1.gotmpl", Output: "output-1.go"},
			{Input: "input-2.gotmpl", Output: "output-2.go"},
		},
	}

	data, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	path := filepath.Join(".forge", "sample.yaml")
	err = os.WriteFile(path, data, 0666)
	if err != nil {
		return err
	}

	fmt.Println("✨ Initialized.")

	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}
