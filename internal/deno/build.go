package deno

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/andymoe/deno-buildpack/internal/metadata"
	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/fs"
)

// Build returns a BuildFunc that provides the deno layer
func Build() packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		config, err := metadata.Read(filepath.Join(context.CNBPath, "buildpack.toml"))
		if err != nil {
			return packit.BuildResult{}, err
		}

		uri := config.Metadata.Dependencies[0].Source

		denoLayer, err := context.Layers.Get("deno", packit.LaunchLayer)
		if err != nil {
			return packit.BuildResult{}, err
		}

		err = denoLayer.Reset()
		if err != nil {
			return packit.BuildResult{}, err
		}

		buildResult, err := InstallDeno(uri, denoLayer)
		if err != nil {
			return buildResult, err
		}

		command := "deno run --allow-all main.ts"
		return packit.BuildResult{
			Plan: context.Plan,
			Layers: []packit.Layer{
				denoLayer,
			},
			Processes: []packit.Process{
				{
					Type:    "web",
					Command: command,
				},
			},
		}, nil
	}
}

// InstallDeno downloads and installs the deno dependency
func InstallDeno(uri string, denoLayer packit.Layer) (packit.BuildResult, error) {
	downloadDir, err := ioutil.TempDir("", "downloadDir")
	if err != nil {
		return packit.BuildResult{}, err
	}
	defer os.RemoveAll(downloadDir)

	fmt.Println("Downloading...")
	fmt.Printf("URI -> %s\n", uri)
	err = exec.Command("curl", "-L", uri,
		"--output", filepath.Join(downloadDir, "deno.gz")).Run()
	if err != nil {
		return packit.BuildResult{}, fmt.Errorf("failed to download deno with error: %w", err)
	}

	fmt.Println("Unziping...")
	err = exec.Command("gunzip", "-d", filepath.Join(downloadDir, "deno.gz")).Run()
	if err != nil {
		return packit.BuildResult{}, fmt.Errorf("Failed to unzip: %w", err)
	}

	err = os.MkdirAll(filepath.Join(denoLayer.Path, "bin"), os.ModePerm)
	if err != nil {
		return packit.BuildResult{}, fmt.Errorf("Failed to make bin dir in deno layer path: %w", err)
	}

	denoExeFile := filepath.Join(denoLayer.Path, "bin", "deno")
	err = fs.Copy(filepath.Join(downloadDir, "deno"), denoExeFile)
	if err != nil {
		return packit.BuildResult{}, fmt.Errorf("Failed moving deno binary to denoLayer path: %w", err)
	}

	err = os.Chmod(denoExeFile, 0550)
	if err != nil {
		return packit.BuildResult{}, fmt.Errorf("Failed to make the deno binary executable: %w", err)
	}
	return packit.BuildResult{}, nil
}
