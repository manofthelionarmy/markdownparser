/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"
	"os"
	"sync"

	"github.com/manofthelionarmy/markdownparser/pkg/parsing/markdown"
	"github.com/manofthelionarmy/markdownparser/pkg/pre"
	"github.com/spf13/cobra"
)

// parseCmd represents the parse command
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		wg := &sync.WaitGroup{}
		trimProcessor := pre.NewTrimProcessor(
			pre.WithWaitGroup(wg),
		)

		cfg := &markdown.Config{}

		markdownFile, err := cmd.Flags().GetString("file")
		if err != nil {
			log.Fatal(err)
		}

		f, err := os.Open(markdownFile)

		markDownParser := markdown.NewMarkdownParser(
			cfg,
			markdown.WithSource(f),
			markdown.WithPreprocessor(trimProcessor),
			markdown.WithTarget(os.Stdout),
			markdown.WithWaitGroup(wg),
		)

		markDownParser.Parse()
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// parseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// parseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	parseCmd.Flags().String("file", "", "Pass in markdown file")
}
