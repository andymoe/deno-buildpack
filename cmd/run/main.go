package main

import (
	"github.com/andymoe/deno-buildpack/internal/deno"
	"github.com/paketo-buildpacks/packit"
)

func main() {
	packit.Run(deno.Detect(), deno.Build())
}
