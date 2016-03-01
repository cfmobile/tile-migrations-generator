package migrations

import (
	"archive/zip"
	"errors"
	"io/ioutil"
	"regexp"

	"gopkg.in/yaml.v2"
)

type ProductVersionFetcher interface {
	FetchProductVersion(productPath string) (string, error)
}

func NewProductVersionFetcher() ProductVersionFetcher {
	return &versionFetcher{}
}

type versionFetcher struct {
}

func (v *versionFetcher) FetchProductVersion(path string) (string, error) {
	m, err := getPackageMetadata(path)
	if err != nil {
		return "", err
	}

	return m.ProductVersion, nil
}

type metadata struct {
	ProductVersion string `yaml:"product_version"`
}

func getPackageMetadata(path string) (*metadata, error) {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}

	var metadataFile *zip.File = nil
	for _, file := range reader.File {
		match, _ := regexp.MatchString(`metadata/.*\.yml`, file.Name)
		if match {
			metadataFile = file
			break
		}
	}

	if metadataFile == nil {
		return nil, errors.New("Could not find package meta data")
	}

	metadataReader, err := metadataFile.Open()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(metadataReader)
	if err != nil {
		return nil, err
	}

	m := metadata{}
	yaml.Unmarshal(data, &m)

	return &m, nil
}
