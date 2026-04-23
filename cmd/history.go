package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Domaa2022/deployctl/internal/history"
)

var historyCmd = &cobra.Command{
	Use:   "history [nombre]",
	Short: "Show deployment history for a container",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		nombre := args[0]

		entries, err := history.Load()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error loading history:", err)
			os.Exit(1)
		}

		// filtrar por nombre
		var matches []history.Entry
		for _, e := range entries {
			if e.Name == nombre {
				matches = append(matches, e)
			}
		}

		if len(matches) == 0 {
			fmt.Printf("No deployment history found for '%s'\n", nombre)
			return
		}

		// imprimir encabezado
		fmt.Printf("Deployment history for '%s'\n\n", nombre)
		fmt.Printf("%-30s %-30s\n", "DATE", "IMAGE")
		fmt.Println("------------------------------------------------------------")

		// imprimir cada entrada — más reciente al final
		for i, e := range matches {
			marker := "  "
			if i == len(matches)-1 {
				marker = "▶ " // marcar el actual
			}
			fmt.Printf("%s%-30s %-30s\n",
				marker,
				e.DeployedAt.Format("2006-01-02 15:04:05"),
				e.Image,
			)
		}
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)
}
