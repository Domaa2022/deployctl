package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/moby/moby/client"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status [nombre]",
	Short: "Show detailed status of a container",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		nombre := args[0]

		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error connecting to Docker:", err)
			os.Exit(1)
		}
		defer cli.Close()

		ctx := context.Background()

		// inspeccionar el contenedor
		info, err := cli.ContainerInspect(ctx, nombre, client.ContainerInspectOptions{})
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		// estado
		status := "stopped"
		if info.Container.State.Running {
			status = "running"
		}

		fmt.Printf("\nContainer: %s\n", nombre)
		fmt.Println("─────────────────────────────────────")
		fmt.Printf("  ID:       %s\n", info.Container.ID[:12])
		fmt.Printf("  Image:    %s\n", info.Container.Config.Image)
		fmt.Printf("  Status:   %s\n", status)
		fmt.Printf("  Started:  %s\n", info.Container.State.StartedAt)
		fmt.Printf("  Restarts: %d\n", info.Container.RestartCount)
		fmt.Println("─────────────────────────────────────")
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
