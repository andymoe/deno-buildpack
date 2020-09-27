package main

import (
	"github.com/andymoe/deno-buildpack/deno"
	"github.com/paketo-buildpacks/packit"
)

func main() {
	packit.Detect(deno.Detect())
}
