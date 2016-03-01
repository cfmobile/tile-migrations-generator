package migrations_test

import (
	"io/ioutil"
	"os"

	. "github.com/cfmobile/tile-migrations-generator/migrations"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExistingMigrations", func() {
	Context("GetExistingMigrations", func() {
		var path string

		It("returns an error if the path is empty", func() {
			path = ""
			migrations, err := GetExistingMigrations(path)
			Expect(err).To(HaveOccurred())
			Expect(migrations).To(BeEmpty())
		})

		It("returns an error if the path doesn't exist", func() {
			path = "/tmp/asdv"
			migrations, err := GetExistingMigrations(path)
			Expect(err).To(HaveOccurred())
			Expect(migrations).To(BeEmpty())
		})

		Context("Path exists", func() {
			BeforeEach(func() {
				var err error
				path, err = ioutil.TempDir("", "")
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				os.RemoveAll(path)
			})

			It("returns an empty list if no migrations exist", func() {
				migrations, err := GetExistingMigrations(path)
				Expect(err).ToNot(HaveOccurred())
				Expect(migrations).To(BeEmpty())
			})

			Context("Existing migrations", func() {
				migrationContent := `  - from_version: 1.1.0.RELEASE
    rules:
    - type: update
      selector: 'product_version'
      to: ((PRODUCT_VERSION))`

				BeforeEach(func() {
					ioutil.WriteFile(path+"/from_1.0.1.RELEASE.yml", []byte(migrationContent), 0666)
				})

				It("Returns the list of migrations", func() {
					migrations, err := GetExistingMigrations(path)
					Expect(err).ToNot(HaveOccurred())
					Expect(migrations).To(ContainElement("1.0.1"))
				})

				It("Returns multiple migrations", func() {
					ioutil.WriteFile(path+"/from_2.0.0.adssda.yml", []byte(migrationContent), 0666)

					migrations, err := GetExistingMigrations(path)
					Expect(err).ToNot(HaveOccurred())
					Expect(migrations).To(ContainElement("1.0.1"))
					Expect(migrations).To(ContainElement("2.0.0"))
				})

				It("Ignores files that don't match the migration format", func() {
					ioutil.WriteFile(path+"/from2.0.0.adssda.yml", []byte(migrationContent), 0666)

					migrations, err := GetExistingMigrations(path)
					Expect(err).ToNot(HaveOccurred())
					Expect(migrations).To(ContainElement("1.0.1"))
					Expect(migrations).To(HaveLen(1))
				})
			})
		})
	})
})
