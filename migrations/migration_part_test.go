package migrations_test

import (
	"io/ioutil"
	"os"
	"time"

	. "github.com/cfmobile/tile-migrations-generator/migrations"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MigrationPart", func() {
	var migrationPart MigrationPart

	BeforeEach(func() {
		migrationPart = MigrationPart{ProductVersion: "some-version"}
	})
	It("Returns an error if dir string is empty", func() {
		err := migrationPart.Write("")
		Expect(err).To(HaveOccurred())
	})
	It("Returns an error if dir does not exist", func() {
		err := migrationPart.Write("/tmp/" + time.Now().String())
		Expect(err).To(HaveOccurred())
	})

	Context("Valid Migration Parts Dir", func() {
		var path string
		BeforeEach(func() {
			var err error
			path, err = ioutil.TempDir("", "")
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			os.RemoveAll(path)
		})

		It("Should create a correctly named migration file", func() {
			migrationPart = MigrationPart{ProductVersion: "test.version.1"}
			err := migrationPart.Write(path)
			Expect(err).ToNot(HaveOccurred())
			migrationFileName := path + "/from_test.version.1.yml"
			Expect(migrationFileName).To(BeARegularFile())
		})
		It("Returns writes a migration from the specific version", func() {
			migrationPart = MigrationPart{ProductVersion: "test.version.2"}
			err := migrationPart.Write(path)
			Expect(err).ToNot(HaveOccurred())

			migrationFile, err := os.Open(path + "/from_test.version.2.yml")
			Expect(err).ToNot(HaveOccurred())
			migrationText, err := ioutil.ReadAll(migrationFile)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(migrationText)).To(Equal(`  - from_version: test.version.2
    rules:
    - type: update
      selector: 'product_version'
      to: ((PRODUCT_VERSION))((BUILD_NUMBER))`))
		})
	})
})
