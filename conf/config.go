package conf

import (
	"crypto/sha512"
	"fmt"
	"io"
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

func loadCheckConfigs(folder string) ([]*cmd.MetaCheck, error) {
	configDir := os.DirFS(folder)

	visitedFileHashes := []string{}
	checks := make([]*cmd.MetaCheck, 0)
	err := fs.WalkDir(configDir, ".", func(path string, dirEntry fs.DirEntry, e error) error {
		// abort if there was an error
		if e != nil {
			return e
		}

		// skip directories and files that are not yaml files
		if dirEntry.IsDir() {
			fmt.Println("Skipping directory: ", path)
			return nil
		} else if !strings.HasSuffix(path, ".yml") && !strings.HasSuffix(path, ".yaml") {
			fmt.Println("Skipping non-yaml file: ", path)
			return nil
		}

		f, err := configDir.Open(path)
		if err != nil {
			return newGenericConfigError(err, path)
		}
		defer f.Close()

		// create hash from file
		hash, err := hashFileSHA512(f.(*os.File))
		if err != nil {
			return newGenericConfigError(err, path)
		}

		// check if the file was already visited
		if slices.Contains(visitedFileHashes, hash) {
			fmt.Println("Skipping duplicate file: ", path)
			return nil
		}
		visitedFileHashes = append(visitedFileHashes, hash)

		fmt.Println("Loading check config: ", path)
		cfg := make(map[string]*cmd.MetaCheck)
		d := yaml.NewDecoder(f)
		d.KnownFields(true)

		err = d.Decode(&cfg)
		if err != nil {
			return newGenericConfigError(err, path)
		}

		// loop over the checks and add them to the list
		for k, v := range cfg {
			if slices.ContainsFunc(checks, func(c *cmd.MetaCheck) bool {
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

func hashFileSHA512(f *os.File) (string, error) {
	defer f.Seek(0, io.SeekStart)

	hash := sha512.New()
	if _, err := io.Copy(hash, f); err != nil {
		return "", err
	}

	hashInBytes := hash.Sum(nil)
	hashString := fmt.Sprintf("%x", hashInBytes)

	return hashString, nil
}
