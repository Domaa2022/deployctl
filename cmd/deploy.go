package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/spf13/cobra"

	"github.com/Domaa2022/deployctl/internal/config"
	"github.com/Domaa2022/deployctl/internal/history"
)

var deployCmd = &cobra.Command{
	Use:   "deploy [nombre] [imagen]",
	Short: "Deploy a Docker container",
	Run: func(cmd *cobra.Command, args []string) {
		var nombre, imagen string

		// leer el flag --env
		env, _ := cmd.Flags().GetString("env")

		if env != "" {
			// modo config: leer de deployctl.yaml
			cfg, err := config.Load()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error loading config:", err)
				os.Exit(1)
			}

			e, err := cfg.GetEnvironment(env)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				os.Exit(1)
			}

			nombre = e.Name
			imagen = e.Image
		} else {
			// modo manual: requiere dos argumentos
			if len(args) < 2 {
				fmt.Fprintln(os.Stderr, "Usage: deployctl deploy [nombre] [imagen]")
				fmt.Fprintln(os.Stderr, "       deployctl deploy --env [dev|staging|prod]")
				os.Exit(1)
			}
			nombre = args[0]
			imagen = args[1]
		}

		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error connecting to Docker:", err)
			os.Exit(1)
		}
		defer cli.Close()

		ctx := context.Background()

		// Paso 1: descargar imagen
		fmt.Printf("⬇  Pulling image %s...\n", imagen)
		reader, err := cli.ImagePull(ctx, imagen, client.ImagePullOptions{})
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error pulling image:", err)
			os.Exit(1)
		}
		io.Copy(io.Discard, reader)
		reader.Close()
		fmt.Println("✓  Image pulled")

		// Paso 2: detener y eliminar contenedor anterior
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

		// Paso 4: guardar en historial
		if err := history.Add(nombre, imagen); err != nil {
			fmt.Fprintln(os.Stderr, "Warning: could not save history:", err)
		}

		fmt.Printf("✓  Container '%s' deployed successfully\n", nombre)
		fmt.Printf("   ID: %s\n", resp.ID[:12])
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().String("env", "", "Environment to deploy (dev|staging|prod)")
}
