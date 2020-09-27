package deno

import (
	"github.com/paketo-buildpacks/packit"
)

// Detect figures out if we are working with a deno app
func Detect() packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
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
}
