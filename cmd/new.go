package cmd

import (
	"fmt"
	"os"

	"github.com/KazukiHayase/forge/codegen"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [generator_name]",
	Short: "Automatically generate files from templates using the forge command",
	Long: `The forge command is a tool for automatically generating files from specified templates. 
With this tool, you can easily generate personalized files by specifying templates from the command line.
Each generator applies different rules for file generation. Make sure to provide the appropriate generator name.
If prompts are configured, you may be prompted to interactively input data during file generation.
For detailed configuration and information on creating generators, please refer to the relevant documentation.`,
	DisableFlagParsing: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Because the flag parsing is disabled, the help flag is not parsed.
		// If the help flag is specified, show the help message and exit.
		for _, a := range args {
			if a == "-h" || a == "--help" {
				cmd.Help()
				os.Exit(0)
			}
		}

		return nil
	},
	RunE: runNewCmd,
}

func runNewCmd(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("generator is required")
	}

	generatorName := args[0]
	generator, err := codegen.NewGenerator(generatorName)
	if err != nil {
		return err
	}

	// If 'flagArgs' is not empty, it is used as the data to omit prompts.
	data := make(codegen.Data)
	flagArgs := args[1:]
	if len(flagArgs) > 0 {
		for _, p := range generator.Config.Prompts {
			cmd.Flags().String(p.Name, "", "")
		}

		cmd.DisableFlagParsing = false
		if err := cmd.ParseFlags(flagArgs); err != nil {
			return err
		}

		for _, p := range generator.Config.Prompts {
			val, err := cmd.Flags().GetString(p.Name)
			if err != nil {
				return err
			}
			if val == "" {
				continue
			}

			data[p.Name] = val
		}
	}

	generator.BuildData(data)

	if err := generator.ParseInOuts(); err != nil {
		return err
	}

	if err := generator.Generate(); err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(newCmd)
}
