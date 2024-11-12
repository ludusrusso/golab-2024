package main

import (
	"golab-2024/pub"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "demo",
	Short: "demo for GoLab 2024",
}

func init() {
	rootCmd.AddCommand(pub.GenerateCMD)
	rootCmd.AddCommand(pub.ProcessCMD)
	rootCmd.AddCommand(pub.StatsCMD)
}

func main() {
	rootCmd.Execute()
}
