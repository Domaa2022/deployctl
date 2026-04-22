package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy [nombre] [imagen]",
	Short: "Deploy a Docker container",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		nombre := args[0]
		imagen := args[1]

		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error connecting to Docker:", err)
			os.Exit(1)
		}
		defer cli.Close()

		ctx := context.Background()

		// Paso 1: descargar imagen
		fmt.Printf("⬇  Pulling image %s...\n", imagen)
		_, err = cli.ImagePull(ctx, imagen, client.ImagePullOptions{})
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error pulling image:", err)
			os.Exit(1)
		}
		fmt.Println("✓  Image pulled")

		// Paso 2: detener y eliminar contenedor anterior si existe
		fmt.Printf("⏹  Stopping existing container '%s'...\n", nombre)
		_, _ = cli.ContainerStop(ctx, nombre, client.ContainerStopOptions{})
		_, _ = cli.ContainerRemove(ctx, nombre, client.ContainerRemoveOptions{})
		fmt.Println("✓  Old container removed")

		// Paso 3: crear y arrancar el nuevo contenedor
		fmt.Printf("🚀 Starting new container '%s'...\n", nombre)
		resp, err := cli.ContainerCreate(ctx,
			client.ContainerCreateOptions{
				Name:   nombre,
				Config: &container.Config{Image: imagen},
			},
		)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating container:", err)
			os.Exit(1)
		}

		_, err = cli.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{})
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error starting container:", err)
			os.Exit(1)
		}

		fmt.Printf("✓  Container '%s' deployed successfully\n", nombre)
		fmt.Printf("   ID: %s\n", resp.ID[:12])
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
