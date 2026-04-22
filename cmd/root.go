package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "deployctl",
	Short: "Deploy and rollback Docker containers without the manual commands",
	Long:  "deployctl automates Docker and AWS ECS deployments with built-in rollback support.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	
}