package main

import (
	_ "github.com/rl404/image-randomizer/docs"
	"github.com/rl404/image-randomizer/internal/utils"
	"github.com/spf13/cobra"
)

// @title Image Randomizer API
// @description Image randomizer API.
// @BasePath /
// @schemes http https
func main() {
	cmd := cobra.Command{
		Use:   "image-randomizer",
		Short: "Image randomizer API",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "server",
		Short: "Run API server",
		RunE: func(*cobra.Command, []string) error {
			return server()
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "migrate",
		Short: "Migrate database",
		RunE: func(*cobra.Command, []string) error {
			return migrate()
		},
	})

	if err := cmd.Execute(); err != nil {
		utils.Fatal(err.Error())
	}
}
