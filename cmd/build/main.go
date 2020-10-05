package main

import (
	"github.com/andymoe/deno-buildpack/internal/deno"
	"github.com/paketo-buildpacks/packit"
)

func main() {
	packit.Build(deno.Build())
}
