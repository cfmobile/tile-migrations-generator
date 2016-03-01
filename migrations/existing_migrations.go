package migrations

import (
	"errors"
	"io/ioutil"
	"regexp"
)

func GetExistingMigrations(path string) ([]string, error) {
	if path == "" {
		return []string{}, errors.New("Path cannot be empty")
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return []string{}, err
	}

	var versions []string
	re := regexp.MustCompile(`from_(\w+\.\w+\.\w+).*\.yml`)
	for _, fileInfo := range files {
		matches := re.FindStringSubmatch(fileInfo.Name())
		if matches != nil {
			versions = append(versions, matches[1])
		}
	}
	return versions, nil
}

func CreateMigrationFromVersion(version string, path string) error {
	return nil
}
