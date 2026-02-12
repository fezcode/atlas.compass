//go:build ignore

package main

import (
	"github.com/fezcode/gobake"
)

func main() {
	bake := gobake.NewEngine()
	bake.LoadRecipeInfo("recipe.piml")

	bake.Task("build", "Builds the binary for multiple platforms", func(ctx *gobake.Context) error {
		ctx.Log("Building %s v%s...", bake.Info.Name, bake.Info.Version)
		
		targets := []struct {
			os   string
			arch string
		}{
			{"linux", "amd64"},
			{"linux", "arm64"},
			{"windows", "amd64"},
			{"windows", "arm64"},
			{"darwin", "amd64"},
			{"darwin", "arm64"},
		}

		err := ctx.Mkdir("build")
		if err != nil {
			return err
		}

		for _, t := range targets {
			output := "build/" + bake.Info.Name + "-" + t.os + "-" + t.arch
			if t.os == "windows" {
				output += ".exe"
			}
			
			ctx.Log("Baking binary for %s/%s -> %s", t.os, t.arch, output)
			
			// Set environment for this specific run
			ctx.Env = []string{
				"GOOS=" + t.os,
				"GOARCH=" + t.arch,
				"CGO_ENABLED=0",
			}
			
			err := ctx.Run("go", "build", "-o", output, "./cmd/main.go")
			if err != nil {
				return err
			}
		}
		return nil
	})

	bake.Task("clean", "Removes build artifacts", func(ctx *gobake.Context) error {
		return ctx.Remove("build")
	})

	bake.Execute()
}
