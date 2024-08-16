package conf

import (
	"io/fs"
	"os"
	"slices"
	"strings"

	"github.com/aimotrens/scrutzone/cmd"
	"gopkg.in/yaml.v3"
)

func Load(configFilePath string) (*cmd.Config, error) {
	mainConfig, err := loadMainConfig(configFilePath)
	if err != nil {
		return nil, err
	}

	mainConfig.Checks, err = loadCheckConfigs(mainConfig.CheckConfigDir)
	if err != nil {
		return nil, err
	}

	return mainConfig, nil
}

func loadMainConfig(configFilePath string) (*cmd.Config, error) {
	f, err := os.OpenFile(configFilePath, os.O_RDONLY, 0)
	if err != nil {
		return nil, newGenericConfigError(err, configFilePath)
	}
	defer f.Close()

	d := yaml.NewDecoder(f)
	d.KnownFields(true)

	config := new(cmd.Config)
	err = d.Decode(config)
	if err != nil {
		return nil, newGenericConfigError(err, configFilePath)
	}

	return config, nil
}

func loadCheckConfigs(folder string) ([]*cmd.Check, error) {
	configDir := os.DirFS(folder)

	checks := make([]*cmd.Check, 0)
	err := fs.WalkDir(configDir, ".", func(path string, dirEntry fs.DirEntry, e error) error {
		// abort if there was an error
		if e != nil {
			return e
		}

		// skip directories and files that are not yaml files
		if dirEntry.IsDir() {
			return nil
		} else if !strings.HasSuffix(path, ".yml") && !strings.HasSuffix(path, ".yaml") {
			return nil
		}

		f, err := configDir.Open(path)
		if err != nil {
			return newGenericConfigError(err, path)
		}
		defer f.Close()

		cfg := make(map[string]*cmd.Check)
		d := yaml.NewDecoder(f)
		d.KnownFields(true)

		err = d.Decode(&cfg)
		if err != nil {
			return newGenericConfigError(err, path)
		}

		for k, v := range cfg {
			if slices.ContainsFunc(checks, func(c *cmd.Check) bool {
				return c.Name == k
			}) {
				return newDuplicateKeyError(k, path)
			}

			// store the key as check name
			v.Name = k

			checks = append(checks, v)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return checks, nil
}
