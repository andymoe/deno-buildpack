package main

import (
	"github.com/andymoe/deno-buildpack/deno"
	"github.com/paketo-buildpacks/packit"
)

func main() {
	packit.Build(deno.Build())
}
