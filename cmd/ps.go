package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/moby/moby/client"
	"github.com/spf13/cobra"
)

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "List running Docker containers",
	Run: func(cmd *cobra.Command, args []string) {
		// Crear cliente de Docker
		cli, err := client.New(client.FromEnv)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error connecting to Docker:", err)
			os.Exit(1)
		}
		defer cli.Close()

		// Listar contenedores corriendo
		result, err := cli.ContainerList(context.Background(), client.ContainerListOptions{})
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error listing containers:", err)
			os.Exit(1)
		}

		// Si no hay contenedores
		if len(result.Items) == 0 {
			fmt.Println("No running containers found.")
			return
		}

		// Imprimir encabezado
		fmt.Printf("%-20s %-30s %-15s\n", "CONTAINER ID", "IMAGE", "STATUS")
		fmt.Println("--------------------------------------------------------------")

		// Imprimir cada contenedor
		for _, c := range result.Items {
			fmt.Printf("%-20s %-30s %-15s\n",
				c.ID[:12],
				c.Image,
				c.Status,
			)
		}
	},
}

func init() {
	rootCmd.AddCommand(psCmd)
}
