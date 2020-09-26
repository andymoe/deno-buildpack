package deno

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/paketo-buildpacks/packit"
)

// Build returns a BuildFunc that provides the deno layer
func Build() packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		file, err := os.Open(filepath.Join(context.CNBPath, "buildpack.toml"))
		if err != nil {
			return packit.BuildResult{}, err
		}

		var m struct {
			Metadata struct {
				Dependencies []struct {
					URI string `toml:"uri"`
				} `toml:"dependencies"`
			} `toml:"metadata"`
		}

		_, err = toml.DecodeReader(file, &m)
		if err != nil {
			return packit.BuildResult{}, err
		}

		uri := m.Metadata.Dependencies[0].URI
		fmt.Printf("URI -> %s\n", uri)

		denoLayer, err := context.Layers.Get("deno", packit.LaunchLayer)
		if err != nil {
			return packit.BuildResult{}, err
		}

		err = denoLayer.Reset()
		if err != nil {
			return packit.BuildResult{}, err
		}

		downloadDir, err := ioutil.TempDir("", "downloadDir")
		if err != nil {
			return packit.BuildResult{}, err
		}
		defer os.RemoveAll(downloadDir)

		fmt.Println("Downloading...")
		err = exec.Command("curl", "-L", uri,
			"--output", filepath.Join(downloadDir, "deno.gz")).Run()
		if err != nil {
			return packit.BuildResult{}, fmt.Errorf("failed to download deno with error: %v", err)
		}

		fmt.Println("unziping...")
		err = exec.Command("gunzip", "-d", filepath.Join(downloadDir, "deno.gz")).Run()
		if err != nil {
			return packit.BuildResult{}, fmt.Errorf("failed to unzip with error: %v", err)
		}

		err = exec.Command("mv", filepath.Join(downloadDir, "deno"), filepath.Join(denoLayer.Path, "deno")).Run()
		if err != nil {
			return packit.BuildResult{}, fmt.Errorf("failed moving deno binary to denoLayer path: %v", err)
		}

		return packit.BuildResult{
			Plan: context.Plan,
			Layers: []packit.Layer{
				denoLayer,
			},
		}, nil
	}
}
