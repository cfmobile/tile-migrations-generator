package migrations

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cfmobile/gopivnet/api"
)

type Migrations struct {
	api            api.Api
	versionFetcher ProductVersionFetcher
	path           string
}

func New(api api.Api, versionFetcher ProductVersionFetcher, path string) (*Migrations, error) {
	if path == "" {
		return nil, errors.New("Migration parts path cannot be empty")
	}
	if api == nil {
		return nil, errors.New("Pivnet api cannot be nil")
	}
	if versionFetcher == nil {
		return nil, errors.New("versionFetcher cannot be nil")
	}
	if _, err := ioutil.ReadDir(path); err != nil {
		return nil, fmt.Errorf(`Provided migration parts path "%s" does not exist`, path)
	}

	return &Migrations{api: api, versionFetcher: versionFetcher, path: path}, nil
}

func (m *Migrations) WriteMissingMigrations(productName string) error {
	missingVersions, err := m.getMissingVersions(productName)
	if err != nil {
		return err
	}

	fmt.Printf("Missing versions %v\n", missingVersions)

	downloadDir, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}

	for _, version := range missingVersions {
		productFile, err := m.api.GetProductFileForVersion(productName, version, "pivotal")
		if err != nil {
			fmt.Printf("Unable to find a product file for %s\n", version)
			continue
		}

		productFilePath := downloadDir + "/" + productFile.Name()

		err = m.api.Download(productFile, productFilePath)
		if err != nil {
			return err
		}
		productVersion, err := m.versionFetcher.FetchProductVersion(productFilePath)
		if err != nil {
			return err
		}

		migrationPart := MigrationPart{ProductVersion: productVersion}
		migrationPart.Write(m.path)

		err = os.Remove(productFilePath)
		if err != nil {
			return err
		}
	}

	os.RemoveAll(downloadDir)

	return nil
}

func (m *Migrations) getMissingVersions(productName string) ([]string, error) {
	existingMigrationVersions, err := GetExistingMigrations(m.path)
	if err != nil {
		return nil, err
	}
	pivnetVersions, err := m.api.GetVersionsForProduct(productName)
	if err != nil {
		return nil, err
	}
	missingMigrationsSet := make(map[string]string)
	for _, version := range pivnetVersions {
		missingMigrationsSet[version] = version
	}

	for _, migration := range existingMigrationVersions {
		delete(missingMigrationsSet, migration)
	}

	missingMigrations := []string{}
	for _, migrationVersion := range missingMigrationsSet {
		missingMigrations = append(missingMigrations, migrationVersion)
	}

	return missingMigrations, nil
}
