package metadata

import (
	"os"

	"github.com/BurntSushi/toml"
)

// BuildpackConfig holds content from buildpack.toml
type BuildpackConfig struct {
	Metadata struct {
		Dependencies []struct {
			ID      string `toml:"id"`
			Version string `toml:"version"`
			Source  string `toml:"source"`
		} `toml:"dependencies"`
	} `toml:"metadata"`
}

// Read metadata from buildpack.toml
func Read(path string) (BuildpackConfig, error) {
	var bc = BuildpackConfig{}
	file, err := os.Open(path)
	if err != nil {
		return bc, err
	}

	_, err = toml.DecodeReader(file, &bc)
	return bc, err
}
