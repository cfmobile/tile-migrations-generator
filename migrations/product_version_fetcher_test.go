package migrations_test

import (
	. "github.com/cfmobile/tile-migrations-generator/migrations"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ProductVersionFetcher", func() {
	var versionFetcher ProductVersionFetcher
	BeforeEach(func() {
		versionFetcher = NewProductVersionFetcher()
	})

	It("returns an error if the path is empty", func() {
		version, err := versionFetcher.FetchProductVersion("")
		Expect(err).To(HaveOccurred())
		Expect(version).To(BeEmpty())
	})

	It("returns an error if the path does not exist", func() {
		version, err := versionFetcher.FetchProductVersion("/tmp/asdvc")
		Expect(err).To(HaveOccurred())
		Expect(version).To(BeEmpty())
	})

	Context("File exists", func() {
		It("returns an error if the metadata does not exist", func() {
			version, err := versionFetcher.FetchProductVersion("./test_files/p-push-notifications-no-metadata.pivotal")
			Expect(err).To(HaveOccurred())
			Expect(version).To(BeEmpty())
		})

		It("returns the version from the metadata", func() {
			version, err := versionFetcher.FetchProductVersion("./test_files/p-push-notifications-1.4.0.114.pivotal")
			Expect(err).NotTo(HaveOccurred())
			Expect(version).To(Equal("1.4.0.114"))
		})
	})
})
