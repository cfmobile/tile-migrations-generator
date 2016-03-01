package migrations

import (
	"errors"
	"fmt"
	"io/ioutil"
)

const (
	migrationTemplate = `  - from_version: %s
    rules:
    - type: update
      selector: 'product_version'
      to: ((PRODUCT_VERSION))((BUILD_NUMBER))`
)

type MigrationPart struct {
	ProductVersion string
}

func (m *MigrationPart) Write(dir string) error {
	if dir == "" {
		return errors.New("No migration part directory provided.")
	}

	if _, err := ioutil.ReadDir(dir); err != nil {
		return fmt.Errorf("Provided migration part directory %s does not exist.", dir)
	}

	migrationFileName := fmt.Sprintf("%s/%s", dir, generateMigrationPartFilename(m.ProductVersion))

	err := ioutil.WriteFile(migrationFileName, []byte(generateMigrationPartContent(m.ProductVersion)), 0666)
	if err != nil {
		return err
	}

	fmt.Printf("Wrote migration file \"%s\" to disk.\n", migrationFileName)
	return nil
}

func generateMigrationPartFilename(version string) string {
	return fmt.Sprintf(`from_%s.yml`, version)
}

func generateMigrationPartContent(version string) string {
	migrationText := fmt.Sprintf(migrationTemplate, version)

	return migrationText
}
