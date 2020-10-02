package deno

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/packit"
)

// Detect figures out if we are working with a deno app
func Detect() packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		ok := DetectDeno(context.WorkingDir)
		if ok {
			return packit.DetectResult{
				Plan: packit.BuildPlan{
					Provides: []packit.BuildPlanProvision{
						{Name: "deno"},
					},
					Requires: []packit.BuildPlanRequirement{
						{Name: "deno"},
					},
				},
			}, nil
		}

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{Name: "deno"},
				},
				Requires: []packit.BuildPlanRequirement{},
			},
		}, nil
	}
}

// DetectDeno scans *.ts or *.js source looking for signs of deno
// It's probably a little too greedy right now
func DetectDeno(workingDir string) (ok bool) {
	if fileExists(filepath.Join(workingDir, "deps.ts")) {
		return true
	}

	denoFound := false
	err := filepath.Walk(workingDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		ext := filepath.Ext(path)
		if !(".ts" == ext || ".js" == ext) {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if bytes.Contains(scanner.Bytes(), []byte("://deno.land/std")) {
				denoFound = true
				return fmt.Errorf("Deno source found in %v", path)
			}

			if bytes.Contains(scanner.Bytes(), []byte("Deno.")) {
				denoFound = true
				return fmt.Errorf("Deno source found in %v", path)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("There was an error scanning the source code: %v", err)
	}

	return denoFound
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
