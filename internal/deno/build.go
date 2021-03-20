package deno

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/andymoe/deno-buildpack/internal/metadata"
	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/pexec"
	"github.com/paketo-buildpacks/packit/vacation"
)

// Build returns a BuildFunc that provides the deno layer
func Build() packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		config, err := metadata.Read(filepath.Join(context.CNBPath, "buildpack.toml"))
		if err != nil {
			return packit.BuildResult{}, err
		}

		uri := config.Metadata.Dependencies[0].Source

		denoLayer, err := context.Layers.Get("deno")
		if err != nil {
			return packit.BuildResult{}, err
		}

		denoLayer, err = denoLayer.Reset()
		if err != nil {
			return packit.BuildResult{}, err
		}

		denoLayer.Launch = true

		buildResult, err := InstallDeno(uri, denoLayer)
		if err != nil {
			return buildResult, err
		}

		denoCache := filepath.Join(denoLayer.Path, "cache")
		err = os.MkdirAll(denoCache, os.ModePerm)
		if err != nil {
			return packit.BuildResult{}, fmt.Errorf("Failed to make deno cache dir in deno layer path: %w", err)
		}

		denoExe := filepath.Join(denoLayer.Path, "bin", "deno")
		deno := pexec.NewExecutable(denoExe)
		os.Setenv("DENO_DIR", denoCache)
		err = deno.Execute(pexec.Execution{
			Args:   []string{"cache", "main.ts"},
			Stdout: os.Stdout,
			Stderr: os.Stderr,
			Dir:    context.WorkingDir,
		})

		if err != nil {
			return packit.BuildResult{}, fmt.Errorf("Failed to cache project dependencies: %w", err)
		}

		command := fmt.Sprintf(`DENO_DIR="%s" deno run --allow-all --cached-only main.ts`, denoCache)
		return packit.BuildResult{
			Plan:   context.Plan,
			Layers: []packit.Layer{denoLayer},
			Launch: packit.LaunchMetadata{
				Processes: []packit.Process{{Type: "web", Command: command}},
				Slices:    []packit.Slice{},
				Labels:    map[string]string{},
			},
		}, nil
	}
}

func InstallDeno(uri string, denoLayer packit.Layer) (packit.BuildResult, error) {
	downloadDir, err := os.MkdirTemp("", "downloadDir")
	if err != nil {
		return packit.BuildResult{}, err
	}
	defer os.RemoveAll(downloadDir)

	fmt.Println("Downloading...")
	fmt.Printf("URI -> %s\n", uri)

	tr := &http.Transport{
		DisableCompression: true,
	}

	c := &http.Client{
		Transport: tr,
	}

	resp, err := c.Get(uri)
	if err != nil {
		return packit.BuildResult{}, fmt.Errorf("failed to download deno with error: %w", err)
	}

	err = os.MkdirAll(filepath.Join(denoLayer.Path, "bin"), os.ModePerm)
	if err != nil {
		return packit.BuildResult{}, fmt.Errorf("Failed to make bin dir in deno layer path: %w", err)
	}

	destination := filepath.Join(denoLayer.Path, "bin")
	archive := vacation.NewArchive(resp.Body)
	archive.Decompress(destination)

	err = os.Chmod(filepath.Join(destination, "deno"), 0550)
	if err != nil {
		return packit.BuildResult{}, fmt.Errorf("Failed to make the deno binary executable: %w", err)
	}

	return packit.BuildResult{}, nil
}
