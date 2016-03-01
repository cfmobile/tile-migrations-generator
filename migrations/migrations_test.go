package migrations_test

import (
	"io/ioutil"
	"os"

	"github.com/cfmobile/gopivnet/api"
	"github.com/cfmobile/gopivnet/resource"
	. "github.com/cfmobile/tile-migrations-generator/migrations"
	"github.com/cfmobile/tile-migrations-generator/migrations/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Migrations", func() {
	var tempDir string
	var versionFetcher *fakes.FakeProductVersionFetcher

	BeforeEach(func() {
		versionFetcher = new(fakes.FakeProductVersionFetcher)
		var err error
		tempDir, err = ioutil.TempDir("", "")
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll(tempDir)
	})

	Context("Initialization", func() {
		It("returns an error if the path is empty", func() {
			migrations, err := New(api.New(""), versionFetcher, "")
			Expect(err).To(HaveOccurred())
			Expect(migrations).To(BeNil())
		})

		It("returns an error if the path does not exist", func() {
			migrations, err := New(api.New(""), versionFetcher, "/tmp/asdx")
			Expect(err).To(HaveOccurred())
			Expect(migrations).To(BeNil())
		})

		It("returns an error if the path is not a directory", func() {
			tempFile, _ := ioutil.TempFile("", "")
			Expect(tempFile).ToNot(BeNil())
			migrations, err := New(api.New(""), versionFetcher, tempFile.Name())
			Expect(err).To(HaveOccurred())
			Expect(migrations).To(BeNil())
		})

		It("returns an error if the api is nil", func() {
			tempDir, _ := ioutil.TempDir("", "")
			migrations, err := New(nil, versionFetcher, tempDir)
			Expect(err).To(HaveOccurred())
			Expect(migrations).To(BeNil())
		})

		It("returns an error if the versionFetcher is nil", func() {
			migrations, err := New(api.New(""), nil, tempDir)
			Expect(err).To(HaveOccurred())
			Expect(migrations).To(BeNil())
		})

		It("returns a Migrations object when passed valid parameters", func() {
			tempDir, _ := ioutil.TempDir("", "")
			migrations, err := New(api.New(""), versionFetcher, tempDir)
			Expect(err).To(BeNil())
			Expect(migrations).ToNot(BeNil())
		})
	})

	Context("Successful initalization", func() {
		var migrations *Migrations
		var api *fakes.FakeApi

		BeforeEach(func() {
			var err error
			existingMigrations := []string{"1.1.1", "1.1.2", "1.1.3"}
			for _, migrationVersion := range existingMigrations {
				fileName := tempDir + "/from_" + migrationVersion + ".extra.stuff.yml"
				ioutil.WriteFile(fileName, []byte{}, 0666)
			}

			api = new(fakes.FakeApi)
			versionFetcher = new(fakes.FakeProductVersionFetcher)

			migrations, err = New(api, versionFetcher, tempDir)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("No missing migrations", func() {
			It("Does not write new migrations if none are returned from pivnet", func() {
				api.GetVersionsForProductReturns([]string{}, nil)

				migrations.WriteMissingMigrations("someProduct")
				files, err := ioutil.ReadDir(tempDir)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(files)).To(Equal(3))
			})

			It("Does not write new migrations if it already has migration parts for them", func() {
				api.GetVersionsForProductReturns([]string{"1.1.1", "1.1.2", "1.1.3"}, nil)

				migrations.WriteMissingMigrations("someProduct")
				files, err := ioutil.ReadDir(tempDir)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(files)).To(Equal(3))
			})
		})

		Context("One missing migration", func() {
			BeforeEach(func() {
				api.GetVersionsForProductReturns([]string{"1.1.1", "1.1.2", "1.1.3", "1.1.4"}, nil)

				api.GetProductFileForVersionReturns(
					&resource.ProductFile{
						AwsObjectKey: "test_files/someProduct.pivotal",
					}, nil)
				versionFetcher.FetchProductVersionReturns("1.1.4_build2", nil)

				api.DownloadStub = func(productFile *resource.ProductFile, path string) error {
					err := ioutil.WriteFile(path, []byte(""), 0666)
					Expect(err).ToNot(HaveOccurred())
					return nil
				}

				migrations.WriteMissingMigrations("someProduct")
			})

			It("Downloads the missing version from pivnet", func() {
				expectedProduct, expectedVersion, _ := api.GetProductFileForVersionArgsForCall(0)
				Expect(expectedProduct).To(Equal("someProduct"))
				Expect(expectedVersion).To(Equal("1.1.4"))

				expectedProductFile, _ := api.DownloadArgsForCall(0)
				Expect(expectedProductFile.Name()).To(Equal("someProduct.pivotal"))
			})

			It("Reads the product version from the product file", func() {
				_, downloadProductPath := api.DownloadArgsForCall(0)
				expectedProductFileName := versionFetcher.FetchProductVersionArgsForCall(0)
				Expect(expectedProductFileName).To(Equal(downloadProductPath))
			})

			It("Writes the migration to disk", func() {
				Expect(tempDir + "/from_1.1.4_build2.yml").To(BeAnExistingFile())
				files, err := ioutil.ReadDir(tempDir)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(files)).To(Equal(4))
			})

			It("Removes the downloaded package file after it finishes", func() {
				_, downloadProductPath := api.DownloadArgsForCall(0)
				Expect(downloadProductPath).ToNot(BeAnExistingFile())
			})
		})
	})
})
