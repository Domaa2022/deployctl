package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/spf13/cobra"

	"github.com/Domaa2022/deployctl/internal/history"
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback [nombre]",
	Short: "Rollback a container to its previous image",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		nombre := args[0]

		// Paso 1: buscar imagen anterior en el historial
		imagen, err := history.Previous(nombre)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		fmt.Printf("⏪  Rolling back '%s' to image %s...\n", nombre, imagen)

		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error connecting to Docker:", err)
			os.Exit(1)
		}
		defer cli.Close()

		ctx := context.Background()

		// Paso 2: descargar imagen anterior
		fmt.Printf("⬇  Pulling image %s...\n", imagen)
		reader, err := cli.ImagePull(ctx, imagen, client.ImagePullOptions{})
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error pulling image:", err)
			os.Exit(1)
		}
		io.Copy(io.Discard, reader)
		reader.Close()
		fmt.Println("✓  Image pulled")

		// Paso 3: detener y eliminar contenedor actual
		fmt.Printf("⏹  Stopping current container '%s'...\n", nombre)
		_, _ = cli.ContainerStop(ctx, nombre, client.ContainerStopOptions{})
		_, _ = cli.ContainerRemove(ctx, nombre, client.ContainerRemoveOptions{})
		fmt.Println("✓  Current container removed")

		// Paso 4: arrancar contenedor con imagen anterior
		fmt.Printf("🚀 Starting container with previous image...\n")
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

		// Paso 5: guardar en historial
		if err := history.Add(nombre, imagen); err != nil {
			fmt.Fprintln(os.Stderr, "Warning: could not save history:", err)
		}

		fmt.Printf("✓  Rollback successful — '%s' is now running %s\n", nombre, imagen)
		fmt.Printf("   ID: %s\n", resp.ID[:12])
	},
}

func init() {
	rootCmd.AddCommand(rollbackCmd)
}
