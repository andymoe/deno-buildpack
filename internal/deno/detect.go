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
		var r []packit.BuildPlanRequirement
		ok := DetectDeno(context.WorkingDir)
		if ok {
			r = append(r, packit.BuildPlanRequirement{
				Name: "deno",
			})
		}

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{Name: "deno"},
				},
				Requires: r,
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

	detected := false
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
				detected = true
				return fmt.Errorf("Deno source found in %v", path)
			}

			if bytes.Contains(scanner.Bytes(), []byte("Deno.")) {
				detected = true
				return fmt.Errorf("Deno source found in %v", path)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("There was an error scanning the source code: %v", err)
	}

	return detected
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
