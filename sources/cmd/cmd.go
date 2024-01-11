package cmd

import (
	"log"
	"os"

	"github.com/ex3ndr/datasets/sync"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func initHandler(cmd *cobra.Command, args []string) error {
	template := []byte("# List your datasets here\ndatasets:\n\n")
	return os.WriteFile("datasets.yaml", template, 0644)
}

func syncHandler(cmd *cobra.Command, args []string) error {

	// Read datasets.yaml
	d, err := os.ReadFile("datasets.yaml")
	if err != nil {
		return err
	}

	// Parse datasets.yaml
	var projectFile sync.ProjectFile
	err = yaml.Unmarshal(d, &projectFile)
	if err != nil {
		return err
	}

	// Sync datasets
	err = sync.Sync(projectFile)
	if err != nil {
		return err
	}

	return nil
}

func NewCLI() *cobra.Command {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cobra.EnableCommandSorting = false

	//
	// Root command
	//

	rootCmd := &cobra.Command{
		Use:           "datasets",
		Short:         "Reproducable datasets for machine learning",
		SilenceUsage:  true,
		SilenceErrors: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Print(cmd.UsageString())
		},
	}

	//
	// Subcommands
	//

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize project",
		RunE:  initHandler,
	}
	rootCmd.AddCommand(initCmd)

	syncCmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync datasets",
		RunE:  syncHandler,
	}
	rootCmd.AddCommand(syncCmd)

	return rootCmd
}
